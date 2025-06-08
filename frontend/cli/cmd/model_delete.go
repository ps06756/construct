package cmd

import (
	"connectrpc.com/connect"
	v1 "github.com/furisto/construct/api/go/v1"
	"github.com/spf13/cobra"
)

func NewModelDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <model-id>",
		Short: "Delete a model by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getAPIClient(cmd.Context())

			_, err := client.Model().DeleteModel(cmd.Context(), &connect.Request[v1.DeleteModelRequest]{
				Msg: &v1.DeleteModelRequest{Id: args[0]},
			})

			return err
		},
	}

	return cmd
}
