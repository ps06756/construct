package cmd

import (
	"fmt"

	"connectrpc.com/connect"
	v1 "github.com/furisto/construct/api/go/v1"
	"github.com/spf13/cobra"
)

func NewMessageDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <message-id>...",
		Short: "Delete one or more messages by ID",
		Long:  `Delete messages by specifying their IDs.`,
		Example: `  # Delete a single message
  construct message delete "123e4567-e89b-12d3-a456-426614174000"`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getAPIClient(cmd.Context())

			for _, messageID := range args {
				_, err := client.Message().DeleteMessage(cmd.Context(), &connect.Request[v1.DeleteMessageRequest]{
					Msg: &v1.DeleteMessageRequest{Id: messageID},
				})

				if err != nil {
					return fmt.Errorf("failed to delete message %s: %w", messageID, err)
				}
			}

			return nil
		},
	}

	return cmd
}
