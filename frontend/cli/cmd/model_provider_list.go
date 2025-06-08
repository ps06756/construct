package cmd

import (
	"connectrpc.com/connect"
	v1 "github.com/furisto/construct/api/go/v1"
	"github.com/spf13/cobra"
)

type modelProviderListOptions struct {
	ShowDisabled  bool
	FormatOptions FormatOptions
}

func NewModelProviderListCmd() *cobra.Command {
	var options modelProviderListOptions

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List model providers",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getAPIClient(cmd.Context())

			filter := &v1.ListModelProvidersRequest_Filter{}
			if !options.ShowDisabled {
				enabled := true
				filter.Enabled = &enabled
			}

			req := &connect.Request[v1.ListModelProvidersRequest]{
				Msg: &v1.ListModelProvidersRequest{
					Filter: filter,
				},
			}

			resp, err := client.ModelProvider().ListModelProviders(cmd.Context(), req)
			if err != nil {
				return err
			}

			displayModelProviders := make([]*ModelProviderDisplay, len(resp.Msg.ModelProviders))
			for i, modelProvider := range resp.Msg.ModelProviders {
				displayModelProviders[i] = ConvertModelProviderToDisplay(modelProvider)
			}

			return getFormatter(cmd.Context()).Display(displayModelProviders, options.FormatOptions.Output)
		},
	}

	cmd.Flags().BoolVar(&options.ShowDisabled, "show-disabled", false, "Show disabled model providers")
	addFormatOptions(cmd, &options.FormatOptions)
	return cmd
}
