package cmd

import (
	"connectrpc.com/connect"
	v1 "github.com/furisto/construct/api/go/v1"
	"github.com/spf13/cobra"
)

type taskGetOptions struct {
	FormatOptions FormatOptions
}

func NewTaskGetCmd() *cobra.Command {
	var options taskGetOptions

	cmd := &cobra.Command{
		Use:   "get <task-id>",
		Short: "Get a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getAPIClient(cmd.Context())

			resp, err := client.Task().GetTask(cmd.Context(), &connect.Request[v1.GetTaskRequest]{
				Msg: &v1.GetTaskRequest{Id: args[0]},
			})

			if err != nil {
				return err
			}

			displayTask := ConvertTaskToDisplay(resp.Msg.Task)
			return getFormatter(cmd.Context()).Display(displayTask, options.FormatOptions.Output)
		},
	}

	addFormatOptions(cmd, &options.FormatOptions)
	return cmd
}
