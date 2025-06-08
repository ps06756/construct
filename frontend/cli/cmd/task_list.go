package cmd

import (
	"connectrpc.com/connect"
	v1 "github.com/furisto/construct/api/go/v1"
	"github.com/spf13/cobra"
)

type taskListOptions struct {
	AgentId       string
	FormatOptions FormatOptions
}

func NewTaskListCmd() *cobra.Command {
	var options taskListOptions

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List tasks",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getAPIClient(cmd.Context())

			filter := &v1.ListTasksRequest_Filter{}
			if options.AgentId != "" {
				filter.AgentId = &options.AgentId
			}

			req := &connect.Request[v1.ListTasksRequest]{
				Msg: &v1.ListTasksRequest{
					Filter: filter,
				},
			}

			resp, err := client.Task().ListTasks(cmd.Context(), req)
			if err != nil {
				return err
			}

			displayTasks := make([]*DisplayTask, len(resp.Msg.Tasks))
			for i, task := range resp.Msg.Tasks {
				displayTasks[i] = ConvertTaskToDisplay(task)
			}

			return getFormatter(cmd.Context()).Display(displayTasks, options.FormatOptions.Output)
		},
	}

	cmd.Flags().StringVarP(&options.AgentId, "agent-id", "a", "", "Filter by agent ID")
	addFormatOptions(cmd, &options.FormatOptions)
	return cmd
}
