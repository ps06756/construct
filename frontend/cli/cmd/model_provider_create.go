package cmd

import (
	"fmt"

	"connectrpc.com/connect"

	v1 "github.com/furisto/construct/api/go/v1"
	"github.com/spf13/cobra"
)

type modelProviderCreateOptions struct {
	Name   string
	ApiKey string
	Type   ModelProviderType
}

func NewModelProviderCreateCmd() *cobra.Command {
	var options modelProviderCreateOptions

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new model provider",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getAPIClient(cmd.Context())

			providerType, err := options.Type.ToAPI()
			if err != nil {
				return err
			}

			resp, err := client.ModelProvider().CreateModelProvider(cmd.Context(), &connect.Request[v1.CreateModelProviderRequest]{
				Msg: &v1.CreateModelProviderRequest{
					Name:           options.Name,
					ProviderType:   providerType,
					Authentication: &v1.CreateModelProviderRequest_ApiKey{ApiKey: options.ApiKey},
				},
			})

			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), resp.Msg.ModelProvider.Id)
			return nil
		},
	}

	cmd.Flags().StringVarP(&options.Name, "name", "n", "", "The name of the model provider")
	cmd.Flags().StringVarP(&options.ApiKey, "api-key", "k", "", "The API key for the model provider")
	cmd.Flags().VarP(&options.Type, "type", "t", "The type of the model provider")

	return cmd
}
