package tool

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/grafana/sobek"
	"github.com/spf13/afero"

	"github.com/furisto/construct/backend/tool/codeact"
)

const editFileDescription = `
## Description
Performs targeted modifications to existing files by replacing specific text sections with new content. This tool enables precise code changes without affecting surrounding content.

## Parameters
- **path** (string, required): Absolute path to the file to modify (e.g., "/workspace/project/src/components/Button.jsx").
- **diffs** (array, required): Array of diff objects, each containing:
  - **old** (string, required): The exact text to find and replace
  - **new** (string, required): The new text to replace it with

## Expected Output
Returns an object indicating success and details about changes made:
%[1]s
{
  "success": true,
  "path": "/path/to/file",
  "replacements_made": 2,
  "expected_replacements": 2
}
%[1]s

## CRITICAL REQUIREMENTS
- **Exact matching**: The "old" content must match file content exactly (whitespace, indentation, line endings)
- **Whitespace preservation**: Maintain proper indentation and formatting in new_text
- **Sufficient context**: Include 3-5 surrounding lines in each "old" text for unique matching
- **Multiple changes**: For multiple changes, add separate objects to the diffs array in file order
- **Concise blocks**: Keep diff blocks focused on specific changes; break large edits into smaller blocks
- **Special operations**:
  - To move code: Use two diffs (one to delete from original (empty "new") + one to insert at new location (empty "old"))
  - To delete code: Use empty string for "new" property
- **File path validation**: Always use absolute paths (starting with "/")

## When to use
- Refactoring code (changing variables, updating functions)
- Bug fixes requiring precise changes
- Feature implementation in existing files
- Configuration changes
- Any targeted code modifications

## Usage Examples

### Single modification
%[1]s
edit_file("/workspace/project/src/utils.js", [
  {
    "old": "function calculateTax(amount) {\n  return amount * 0.08;\n}",
    "new": "function calculateTax(amount, rate = 0.08) {\n  return amount * rate;\n}"
  }
]);
%[1]s

### Multiple modifications
%[1]s
edit_file("/workspace/project/src/components/Button.jsx", [
  {
    "old": "import React from 'react';",
    "new": "import React, { useState } from 'react';"
  },
  {
    "old": "function Button({ text, onClick }) {",
    "new": "function Button({ text, onClick, disabled = false }) {"
  },
  {
    "old": "<button className=\"primary-button\" onClick={onClick}>",
    "new": "<button className=\"primary-button\" onClick={onClick} disabled={disabled}>"
  }
]);
%[1]s
`

func NewEditFileTool() codeact.Tool {
	return codeact.NewOnDemandTool(
		"edit_file",
		fmt.Sprintf(editFileDescription, "```"),
		editFileHandler,
	)
}

type EditFileInput struct {
	Path  string     `json:"path"`
	Diffs []DiffPair `json:"diffs"`
}

type DiffPair struct {
	Old string `json:"old"`
	New string `json:"new"`
}

type DiffValidationError struct {
	DiffIndex    int    `json:"diff_index"`
	ErrorType    string `json:"error_type"`
	ErrorMessage string `json:"error_message"`
	SuggestedFix string `json:"suggested_fix,omitempty"`
}

type ConflictWarning struct {
	Edit1Index   int    `json:"edit1_index"`
	Edit2Index   int    `json:"edit2_index"`
	ConflictType string `json:"conflict_type"`
	Message      string `json:"message"`
}

type EditFileResult struct {
	Success              bool                  `json:"success"`
	Path                 string                `json:"path"`
	ReplacementsMade     int                   `json:"replacements_made"`
	ExpectedReplacements int                   `json:"expected_replacements"`
	FailureReason        string                `json:"failure_reason,omitempty"`
	ValidationErrors     []DiffValidationError `json:"validation_errors,omitempty"`
	ConflictWarnings     []ConflictWarning     `json:"conflict_warnings,omitempty"`
}

func editFileHandler(session *codeact.Session) func(call sobek.FunctionCall) sobek.Value {
	return func(call sobek.FunctionCall) sobek.Value {
		if len(call.Arguments) != 2 {
			session.Throw(codeact.NewCustomError("edit_file requires exactly 2 arguments: path and diffs", []string{
				"- **path** (string, required): Absolute path to the file to modify (e.g., \"/workspace/project/src/components/Button.jsx\").",
				"- **diffs** (array, required): Array of diff objects, each containing: - **old** (string, required): The exact text to find and replace - **new** (string, required): The new text to replace it with",
			}))
		}

		path := call.Argument(0).String()
		diffsArg := call.Argument(1)

		// Parse diffs array
		var diffs []DiffPair
		if diffsObj := diffsArg.ToObject(session.VM); diffsObj != nil {
			if lengthVal := diffsObj.Get("length"); lengthVal != nil {
				length := int(lengthVal.ToInteger())
				for i := 0; i < length; i++ {
					if diffVal := diffsObj.Get(fmt.Sprintf("%d", i)); diffVal != nil {
						if diffObj := diffVal.ToObject(session.VM); diffObj != nil {
							oldText := ""
							newText := ""
							if oldVal := diffObj.Get("old"); oldVal != nil {
								oldText = oldVal.String()
							}
							if newVal := diffObj.Get("new"); newVal != nil {
								newText = newVal.String()
							}
							diffs = append(diffs, DiffPair{Old: oldText, New: newText})
						}
					}
				}
			}
		}

		if len(diffs) == 0 {
			session.Throw(codeact.NewCustomError("diffs array cannot be empty", []string{
				"Provide at least one diff object with 'old' and 'new' properties",
			}))
		}

		result, err := editFile(session.FS, &EditFileInput{
			Path:  path,
			Diffs: diffs,
		})
		if err != nil {
			session.Throw(err)
		}

		return session.VM.ToValue(result)
	}
}

func processEdits(fileContent string, diffs []DiffPair) (string, int, []DiffValidationError) {
	var validationErrors []DiffValidationError
	workingContent := fileContent
	replacementsMade := 0

	for i, diff := range diffs {
		if diff.Old == "" && diff.New == "" {
			continue
		}

		if diff.Old == diff.New {
			validationErrors = append(validationErrors, DiffValidationError{
				DiffIndex:    i,
				ErrorType:    "no_op",
				ErrorMessage: "old and new text are identical",
				SuggestedFix: "Remove this diff or provide different new text",
			})

			continue
		}

		if diff.Old == "" {
			workingContent = workingContent + diff.New
			replacementsMade++
			continue
		}

		if strings.Contains(workingContent, diff.Old) {
			// Exact match found
			if diff.New == "" {
				// Delete operation
				workingContent = strings.Replace(workingContent, diff.Old, "", 1)
			} else {
				// Replace operation
				workingContent = strings.Replace(workingContent, diff.Old, diff.New, 1)
			}
			replacementsMade++

		} else {
			startIdx, endIdx := lineTrimmedFallbackMatch(workingContent, diff.Old, 0)
			if startIdx != -1 && endIdx != -1 {
				// Fallback match found
				if diff.New == "" {
					// Delete operation
					workingContent = workingContent[:startIdx] + workingContent[endIdx:]
				} else {
					// Replace operation
					workingContent = workingContent[:startIdx] + diff.New + workingContent[endIdx:]
				}
				replacementsMade++
			} else {
				validationErrors = append(validationErrors, DiffValidationError{
					DiffIndex:    i,
					ErrorType:    "not_found",
					ErrorMessage: "old text not found in file",
					SuggestedFix: "Check that the old text exactly matches the file content, including whitespace and indentation",
				})
			}
		}
	}

	return workingContent, replacementsMade, validationErrors
}

func lineTrimmedFallbackMatch(originalContent, searchContent string, startIndex int) (int, int) {
	originalLines := strings.Split(originalContent, "\n")
	searchLines := strings.Split(searchContent, "\n")

	// Trim trailing empty line if exists (from the trailing \n in searchContent)
	if len(searchLines) > 0 && searchLines[len(searchLines)-1] == "" {
		searchLines = searchLines[:len(searchLines)-1]
	}

	if len(searchLines) == 0 {
		return -1, -1
	}

	// Find the line number where startIndex falls
	startLineNum := 0
	currentIndex := 0
	for currentIndex < startIndex && startLineNum < len(originalLines) {
		currentIndex += len(originalLines[startLineNum]) + 1 // +1 for \n
		startLineNum++
	}

	// For each possible starting position in original content
	for i := startLineNum; i <= len(originalLines)-len(searchLines); i++ {
		matches := true

		// Try to match all search lines from this position
		for j := 0; j < len(searchLines); j++ {
			originalTrimmed := strings.TrimSpace(originalLines[i+j])
			searchTrimmed := strings.TrimSpace(searchLines[j])

			if originalTrimmed != searchTrimmed {
				matches = false
				break
			}
		}

		// If we found a match, calculate the exact character positions
		if matches {
			// Find start character index
			matchStartIndex := 0
			for k := 0; k < i; k++ {
				matchStartIndex += len(originalLines[k]) + 1 // +1 for \n
			}

			// Find end character index
			matchEndIndex := matchStartIndex
			for k := 0; k < len(searchLines); k++ {
				matchEndIndex += len(originalLines[i+k]) + 1 // +1 for \n
			}

			return matchStartIndex, matchEndIndex
		}
	}

	return -1, -1
}

// detectConflicts analyzes potential conflicts between multiple edits
func detectConflicts(diffs []DiffPair) []ConflictWarning {
	var warnings []ConflictWarning

	for i := 0; i < len(diffs)-1; i++ {
		for j := i + 1; j < len(diffs); j++ {
			edit1 := diffs[i]
			edit2 := diffs[j]

			// Skip empty diffs
			if (edit1.Old == "" && edit1.New == "") || (edit2.Old == "" && edit2.New == "") {
				continue
			}

			// Check if edit j depends on the result of edit i
			if edit2.Old != "" && edit1.New != "" && strings.Contains(edit2.Old, edit1.New) {
				warnings = append(warnings, ConflictWarning{
					Edit1Index:   i + 1,
					Edit2Index:   j + 1,
					ConflictType: "dependency",
					Message:      fmt.Sprintf("Edit %d depends on result of edit %d", j+1, i+1),
				})
			}

			// Check if both edits try to modify overlapping text regions
			if edit1.Old != "" && edit2.Old != "" {
				// Simple overlap detection: check if one old text contains the other
				if strings.Contains(edit1.Old, edit2.Old) || strings.Contains(edit2.Old, edit1.Old) {
					warnings = append(warnings, ConflictWarning{
						Edit1Index:   i + 1,
						Edit2Index:   j + 1,
						ConflictType: "overlap",
						Message:      fmt.Sprintf("Edits %d and %d may affect overlapping text regions", i+1, j+1),
					})
				}

				// Check if both edits target the same exact text
				if edit1.Old == edit2.Old {
					warnings = append(warnings, ConflictWarning{
						Edit1Index:   i + 1,
						Edit2Index:   j + 1,
						ConflictType: "duplicate_target",
						Message:      fmt.Sprintf("Edits %d and %d target the same text", i+1, j+1),
					})
				}
			}

			// Check for potential line-level conflicts by examining line boundaries
			if edit1.Old != "" && edit2.Old != "" {
				edit1Lines := strings.Split(edit1.Old, "\n")
				edit2Lines := strings.Split(edit2.Old, "\n")

				// If both edits span multiple lines and share common line content
				if len(edit1Lines) > 1 && len(edit2Lines) > 1 {
					for _, line1 := range edit1Lines {
						for _, line2 := range edit2Lines {
							if strings.TrimSpace(line1) != "" && strings.TrimSpace(line1) == strings.TrimSpace(line2) {
								warnings = append(warnings, ConflictWarning{
									Edit1Index:   i + 1,
									Edit2Index:   j + 1,
									ConflictType: "line_overlap",
									Message:      fmt.Sprintf("Edits %d and %d may affect the same line", i+1, j+1),
								})
								break
							}
						}
					}
				}
			}
		}
	}

	return warnings
}

func editFile(fsys afero.Fs, input *EditFileInput) (*EditFileResult, error) {
	if !filepath.IsAbs(input.Path) {
		return nil, codeact.NewError(codeact.PathIsNotAbsolute, "path", input.Path)
	}
	path := input.Path

	// Check if file exists and is not a directory
	stat, err := fsys.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, codeact.NewError(codeact.FileNotFound, "path", path)
		}
		return nil, codeact.NewCustomError("error accessing file", []string{
			"Verify that you have the permission to access the file",
		}, "path", path, "error", err)
	}

	if stat.IsDir() {
		return nil, codeact.NewCustomError("path is a directory", []string{
			"Please provide a valid path to a file",
		}, "path", path)
	}

	// Read file content
	content, err := afero.ReadFile(fsys, path)
	if err != nil {
		return nil, codeact.NewCustomError("error reading file", []string{
			"Verify that you have the permission to read the file",
		}, "path", path, "error", err)
	}

	fileContent := string(content)
	expectedReplacements := len(input.Diffs)
	
	conflictWarnings := detectConflicts(input.Diffs)

	newContent, replacementsMade, validationErrors := processEdits(fileContent, input.Diffs)
	if len(validationErrors) > 0 {
		return &EditFileResult{
			Success:              false,
			Path:                 path,
			ReplacementsMade:     0,
			ExpectedReplacements: expectedReplacements,
			FailureReason:        fmt.Sprintf("validation failed: %d error(s) found", len(validationErrors)),
			ValidationErrors:     validationErrors,
		}, nil
	}

	// Only write if content actually changed
	if newContent != fileContent {
		err = afero.WriteFile(fsys, path, []byte(newContent), stat.Mode())
		if err != nil {
			return nil, codeact.NewCustomError("error writing file", []string{
				"Verify that you have the permission to write to the file",
			}, "path", path, "error", err)
		}
	}

	return &EditFileResult{
		Success:              true,
		Path:                 path,
		ReplacementsMade:     replacementsMade,
		ExpectedReplacements: expectedReplacements,
		ConflictWarnings:     conflictWarnings,
	}, nil
}
