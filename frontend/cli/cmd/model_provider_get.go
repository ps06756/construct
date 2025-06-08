package cmd

import (
	"connectrpc.com/connect"
	v1 "github.com/furisto/construct/api/go/v1"
	"github.com/spf13/cobra"
)

type modelProviderGetOptions struct {
	FormatOptions FormatOptions
}

func NewModelProviderGetCmd() *cobra.Command {
	var options modelProviderGetOptions

	cmd := &cobra.Command{
		Use:   "get <model-provider-id>",
		Short: "Get a model provider",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getAPIClient(cmd.Context())

			resp, err := client.ModelProvider().GetModelProvider(cmd.Context(), &connect.Request[v1.GetModelProviderRequest]{
				Msg: &v1.GetModelProviderRequest{Id: args[0]},
			})

			if err != nil {
				return err
			}

			displayModelProvider := ConvertModelProviderToDisplay(resp.Msg.ModelProvider)
			return getFormatter(cmd.Context()).Display(displayModelProvider, options.FormatOptions.Output)
		},
	}

	addFormatOptions(cmd, &options.FormatOptions)
	return cmd
}
