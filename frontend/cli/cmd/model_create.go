package cmd

import (
	"fmt"

	"connectrpc.com/connect"

	v1 "github.com/furisto/construct/api/go/v1"
	"github.com/spf13/cobra"
)

type modelCreateOptions struct {
	Name            string
	ModelProviderID string
	ContextWindow   int64
}

func NewModelCreateCmd() *cobra.Command {
	var options modelCreateOptions

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new model",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getAPIClient(cmd.Context())

			resp, err := client.Model().CreateModel(cmd.Context(), &connect.Request[v1.CreateModelRequest]{
				Msg: &v1.CreateModelRequest{
					Name:            options.Name,
					ModelProviderId: options.ModelProviderID,
					ContextWindow:   options.ContextWindow,
				},
			})

			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), resp.Msg.Model.Id)
			return nil
		},
	}

	cmd.Flags().StringVarP(&options.Name, "name", "n", "", "The name of the model")
	cmd.Flags().StringVarP(&options.ModelProviderID, "model-provider-id", "p", "", "The ID of the model provider")
	cmd.Flags().Int64VarP(&options.ContextWindow, "context-window", "w", 0, "The context window size")

	return cmd
}
