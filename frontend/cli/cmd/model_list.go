package cmd

import (
	"connectrpc.com/connect"

	v1 "github.com/furisto/construct/api/go/v1"
	"github.com/spf13/cobra"
)

type modelListOptions struct {
	ModelProviderID string
	Enabled         bool
	ShowDisabled    bool
	FormatOptions   FormatOptions
}

func NewModelListCmd() *cobra.Command {
	var options modelListOptions

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List models",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getAPIClient(cmd.Context())

			filter := &v1.ListModelsRequest_Filter{}
			if options.ModelProviderID != "" {
				filter.ModelProviderId = &options.ModelProviderID
			}

			if !options.ShowDisabled {
				enabled := true
				filter.Enabled = &enabled
			}

			req := &connect.Request[v1.ListModelsRequest]{
				Msg: &v1.ListModelsRequest{
					Filter: filter,
				},
			}

			resp, err := client.Model().ListModels(cmd.Context(), req)
			if err != nil {
				return err
			}

			displayModels := make([]*ModelDisplay, len(resp.Msg.Models))
			for i, model := range resp.Msg.Models {
				displayModels[i] = ConvertModelToDisplay(model)
			}

			return getFormatter(cmd.Context()).Display(displayModels, options.FormatOptions.Output)
		},
	}

	cmd.Flags().StringVarP(&options.ModelProviderID, "model-provider-id", "p", "", "Filter by model provider ID")
	cmd.Flags().BoolVar(&options.ShowDisabled, "show-disabled", false, "Show disabled models")
	addFormatOptions(cmd, &options.FormatOptions)
	return cmd
}
