package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewResumeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "resume <task-id>",
		Short:   "Resume a task",
		Long:    `The "resume" command allows you to resume a task`,
		GroupID: "core",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Resuming task", args[0])
		},
	}

	return cmd
}
