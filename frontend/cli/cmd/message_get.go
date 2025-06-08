package cmd

import (
	"connectrpc.com/connect"
	v1 "github.com/furisto/construct/api/go/v1"
	"github.com/spf13/cobra"
)

type messageGetOptions struct {
	FormatOptions FormatOptions
}

func NewMessageGetCmd() *cobra.Command {
	var options messageGetOptions

	cmd := &cobra.Command{
		Use:   "get <message-id>",
		Short: "Get a message",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getAPIClient(cmd.Context())

			resp, err := client.Message().GetMessage(cmd.Context(), &connect.Request[v1.GetMessageRequest]{
				Msg: &v1.GetMessageRequest{Id: args[0]},
			})

			if err != nil {
				return err
			}

			displayMessage := ConvertMessageToDisplay(resp.Msg.Message)
			return getFormatter(cmd.Context()).Display(displayMessage, options.FormatOptions.Output)
		},
	}

	addFormatOptions(cmd, &options.FormatOptions)
	return cmd
}
