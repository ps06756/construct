package cmd

import (
	"connectrpc.com/connect"
	v1 "github.com/furisto/construct/api/go/v1"
	"github.com/spf13/cobra"
)

func NewTaskDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <task-id>",
		Short: "Delete a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getAPIClient(cmd.Context())

			_, err := client.Task().DeleteTask(cmd.Context(), &connect.Request[v1.DeleteTaskRequest]{
				Msg: &v1.DeleteTaskRequest{Id: args[0]},
			})

			return err
		},
	}

	return cmd
}
