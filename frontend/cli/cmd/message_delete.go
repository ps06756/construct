package cmd

import (
	"connectrpc.com/connect"
	v1 "github.com/furisto/construct/api/go/v1"
	"github.com/spf13/cobra"
)

func NewMessageDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <message-id>",
		Short: "Delete a message",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getAPIClient(cmd.Context())

			_, err := client.Message().DeleteMessage(cmd.Context(), &connect.Request[v1.DeleteMessageRequest]{
				Msg: &v1.DeleteMessageRequest{Id: args[0]},
			})

			return err
		},
	}

	return cmd
}
