package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewResumeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resume [task-id] [flags]",
		Short: "Resume an existing conversation",
		Long: `Resume an existing conversation from where you left off. If no task ID is provided, shows an interactive picker of recent tasks. Supports partial ID matching for convenience.

The resumed conversation will restore the full context including:
- All previous messages in the conversation
- The agent that was being used
- The workspace directory settings

Examples:
  # Show interactive picker to select from recent tasks
  construct resume

  # Resume the most recent task immediately
  construct resume --last

  # Resume specific task by full ID
  construct resume 01974c1d-0be8-70e1-88b4-ad9462fff25e`,
		GroupID: "core",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Resuming task", args[0])
		},
	}

	cmd.Flags().Bool("last", false, "Resume the most recent task immediately")

	return cmd
}
