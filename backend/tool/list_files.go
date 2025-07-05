package tool

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/furisto/construct/backend/tool/codeact"
	"github.com/grafana/sobek"
	"github.com/spf13/afero"
)

const listFilesDescription = `
## Description
Lists the contents (files and subdirectories) of a specified directory. This tool helps explore project file structures and navigate directories by providing a clear, structured view of their contents.

## Parameters
- **path** (string, required): Absolute path to the directory you want to list (e.g., "/workspace/project/src"). Forward slashes (/) work on all platforms.
- **recursive** (boolean, required): When set to true, lists all files and directories recursively through all subdirectories. When false, only lists the top-level contents of the specified directory.

## Expected Output
Returns an object containing an array of directory entries. A file is identified by the type code "f" and a directory by the type code "d":
%[1]s
{
  "path": "/absolute/path/to/listed/directory", 
  "entries": [
    {
      "n": "/absolute/path/to/listed/directory/file.js", 
      "t": "f", 
      "s": 8 
    },
    {
      "n": "/absolute/path/to/listed/directory/images", 
      "t": "d",
      "s": 0 
    },
    {
      "n": "/absolute/path/to/listed/directory/images/logo.png", 
      "t": "f", 
      "s": 5 
    }
  ]
}
%[1]s

**Details:**
-   The %[2]spath%[2]s field in the response object will be the same absolute path provided in the %[2]spath%[2]s parameter.
-   %[2]sentries%[2]s: An array of objects, each representing a file or directory.
    -   %[2]sn%[2]s (name): The name of the file or subdirectory. This will always be an absolute path.
    -   %[2]st%[2]s (type): A character code indicating the entry type:
        -   %[2]s'f'%[2]s: Represents a regular file.
        -   %[2]s'd'%[2]s: Represents a directory.
    -   %[2]ss%[2]s (size): The size of the entry **in kilobytes**. For directories, the size is reported as 0.
-   If the target directory is empty, %[2]sentries%[2]s will be an empty array (%[2]s[]%[2]s).
-   If the specified %[2]spath%[2]s does not exist, is not a directory, or cannot be accessed due to permissions or other issues, the tool will throw an exception with a descriptive error message.

## IMPORTANT USAGE NOTES
- **Path format**: Always use absolute paths starting with "/"
%[1]s
  // Correct path format
  list_files("/workspace/project/src", false)
%[1]s
- **Performance considerations**: Be cautious with the recursive option on large directories
%[1]s
  // First list non-recursively to understand structure
  try {
    const topLevelContents = list_files("/workspace/project", false);
    print("Top-level directories:", topLevelContents.entries
      .filter(entry => entry.t === "d")
      .map(dir => dir.n));

    // Then list specific subdirectories recursively if needed
    const componentsContents = list_files("/workspace/project/src/components", true);
  } catch (error) {
    print("Error exploring project structure:", error);
  }
%[1]s
- **Exception handling**: Always wrap directory operations in try/catch blocks

## When to use
- **Project exploration**: When you need to understand the structure of a project
- **File location**: When looking for specific files or file types (e.g., all *.js files in /src)
- **Verification**: To confirm directories exist before performing operations
- **Path discovery**: To identify the correct paths for subsequent file operations
- **Pre-computation for operations**: Before batch operations like deleting, copying, or archiving, to gather the list of items to be processed

## Usage Examples

%[1]s
try {
  // List top-level contents non-recursively
  const srcFiles = list_files("/workspace/project/src", false);
  print("Top-level JS files:", srcFiles.entries
    .filter(e => e.t === "f" && e.n.endsWith(".js"))
    .map(f => f.n));

  // Find subdirectories and explore one recursively
  const components = srcFiles.entries.find(e => e.t === "d" && e.n === "components");
  if (components) {
    const allComponents = list_files("/workspace/project/src/components", true);

    // Group files by extension
    const byExt = allComponents.entries
      .filter(e => e.t === "f")
      .reduce((acc, f) => {
        const ext = f.n.split('.').pop() || "unknown";
        acc[ext] = (acc[ext] || 0) + 1;
        return acc;
      }, {});
    print("Files by extension:", byExt);
  }
} catch (error) {
  print("Error listing directory:", error);
}
%[1]s
`

func NewListFilesTool() codeact.Tool {
	return codeact.NewOnDemandTool(
		ToolNameListFiles,
		fmt.Sprintf(listFilesDescription, "```", "`"),
		listFilesHandler,
	)
}

type ListFilesInput struct {
	Path      string
	Recursive bool
}

type ListFilesResult struct {
	Path    string           `json:"path"`
	Entries []DirectoryEntry `json:"entries"`
}

type DirectoryEntry struct {
	Name string `json:"n"`
	Type string `json:"t"`
	Size int64  `json:"s"`
}

func listFilesHandler(session *codeact.Session) func(call sobek.FunctionCall) sobek.Value {
	return func(call sobek.FunctionCall) sobek.Value {
		if len(call.Arguments) != 2 {
			session.Throw(codeact.NewCustomError("list_files requires exactly 2 arguments: path and recursive", []string{
				"- **path** (string, required): Absolute path to the directory you want to list (e.g., \"/workspace/project/src\"). Forward slashes (/) work on all platforms.\n" +
					"- **recursive** (boolean, required): When set to true, lists all files and directories recursively through all subdirectories. When false only lists the top-level contents of the specified directory.",
			}))
		}

		path := call.Argument(0).String()
		recursive := call.Argument(1).ToBoolean()

		result, err := listFiles(session.FS, &ListFilesInput{
			Path:      path,
			Recursive: recursive,
		})
		if err != nil {
			session.Throw(err)
		}

		return session.VM.ToValue(result)
	}
}

func listFiles(fsys afero.Fs, input *ListFilesInput) (*ListFilesResult, error) {
	if !filepath.IsAbs(input.Path) {
		return nil, codeact.NewError(codeact.PathIsNotAbsolute, "path", input.Path)
	}
	path := input.Path

	fileInfo, err := fsys.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, codeact.NewError(codeact.DirectoryNotFound, "path", path)
		}
		if os.IsPermission(err) {
			return nil, codeact.NewError(codeact.PermissionDenied, "path", path)
		}
		return nil, codeact.NewError(codeact.CannotStatFile, "path", path)
	}

	if !fileInfo.IsDir() {
		return nil, codeact.NewError(codeact.PathIsNotDirectory, "path", path)
	}

	entries := []DirectoryEntry{}
	if input.Recursive {
		err = afero.Walk(fsys, path, func(filePath string, entry fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if filePath == path {
				return nil
			}

			dirEntry, err := toDirectoryEntry(filePath, entry)
			if err != nil {
				return err
			}
			entries = append(entries, *dirEntry)
			return nil
		})

		if err != nil {
			if os.IsPermission(err) {
				return nil, codeact.NewError(codeact.PermissionDenied, "path", path)
			}
			return nil, codeact.NewError(codeact.GenericFileError, "path", path, "error", err)
		}
	} else {
		dirEntries, err := afero.ReadDir(fsys, path)
		if err != nil {
			if os.IsPermission(err) {
				return nil, codeact.NewError(codeact.PermissionDenied, "path", path)
			}
			return nil, codeact.NewError(codeact.GenericFileError, "path", path, "error", err)
		}

		for _, entry := range dirEntries {
			entryPath := filepath.Join(path, entry.Name())
			dirEntry, err := toDirectoryEntry(entryPath, entry)
			if err != nil {
				return nil, err
			}
			entries = append(entries, *dirEntry)
		}
	}

	return &ListFilesResult{
		Path:    path,
		Entries: entries,
	}, nil
}

func toDirectoryEntry(path string, info fs.FileInfo) (*DirectoryEntry, error) {
	if info.IsDir() {
		return &DirectoryEntry{
			Name: path,
			Type: "d",
			Size: 0,
		}, nil
	} else {
		return &DirectoryEntry{
			Name: path,
			Type: "f",
			Size: (info.Size() + 1023) / 1024, // Size in KB
		}, nil
	}
}
