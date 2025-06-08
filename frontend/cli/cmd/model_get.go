package cmd

import (
	"connectrpc.com/connect"
	v1 "github.com/furisto/construct/api/go/v1"
	"github.com/spf13/cobra"
)

type modelGetOptions struct {
	FormatOptions FormatOptions
}

func NewModelGetCmd() *cobra.Command {
	var options modelGetOptions

	cmd := &cobra.Command{
		Use:   "get <model-id>",
		Short: "Get a model by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getAPIClient(cmd.Context())

			req := &connect.Request[v1.GetModelRequest]{
				Msg: &v1.GetModelRequest{Id: args[0]},
			}

			resp, err := client.Model().GetModel(cmd.Context(), req)
			if err != nil {
				return err
			}

			displayModel := ConvertModelToDisplay(resp.Msg.Model)
			return getFormatter(cmd.Context()).Display(displayModel, options.FormatOptions.Output)
		},
	}

	addFormatOptions(cmd, &options.FormatOptions)
	return cmd
}
