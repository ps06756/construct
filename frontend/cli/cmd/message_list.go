package cmd

import (
	"connectrpc.com/connect"
	v1 "github.com/furisto/construct/api/go/v1"
	"github.com/spf13/cobra"
)

type messageListOptions struct {
	TaskId        string
	AgentId       string
	Role          string
	FormatOptions FormatOptions
}

func NewMessageListCmd() *cobra.Command {
	var options messageListOptions

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List messages",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getAPIClient(cmd.Context())

			filter := &v1.ListMessagesRequest_Filter{}

			if options.TaskId != "" {
				filter.TaskId = &options.TaskId
			}

			if options.AgentId != "" {
				filter.AgentId = &options.AgentId
			}

			if options.Role != "" {
				switch options.Role {
				case "user":
					role := v1.MessageRole_MESSAGE_ROLE_USER
					filter.Role = &role
				case "assistant":
					role := v1.MessageRole_MESSAGE_ROLE_ASSISTANT
					filter.Role = &role
				}
			}

			req := &connect.Request[v1.ListMessagesRequest]{
				Msg: &v1.ListMessagesRequest{
					Filter: filter,
				},
			}

			resp, err := client.Message().ListMessages(cmd.Context(), req)
			if err != nil {
				return err
			}

			displayMessages := make([]*DisplayMessage, len(resp.Msg.Messages))
			for i, message := range resp.Msg.Messages {
				displayMessages[i] = ConvertMessageToDisplay(message)
			}

			return getFormatter(cmd.Context()).Display(displayMessages, options.FormatOptions.Output)
		},
	}

	cmd.Flags().StringVarP(&options.TaskId, "task-id", "t", "", "Filter by task ID")
	cmd.Flags().StringVarP(&options.AgentId, "agent-id", "a", "", "Filter by agent ID")
	cmd.Flags().StringVarP(&options.Role, "role", "r", "", "Filter by role (user or assistant)")
	addFormatOptions(cmd, &options.FormatOptions)
	return cmd
}
