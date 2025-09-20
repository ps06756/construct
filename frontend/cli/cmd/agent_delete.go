package cmd

import (
	"fmt"

	"connectrpc.com/connect"
	v1 "github.com/furisto/construct/api/go/v1"
	"github.com/spf13/cobra"
)

type agentDeleteOptions struct {
	Force bool
}

func NewAgentDeleteCmd() *cobra.Command {
	options := new(agentDeleteOptions)
	cmd := &cobra.Command{
		Use:     "delete <name|id>... [flags]",
		Short:   "Permanently delete one or more agents",
		Aliases: []string{"rm"},
		Example: `  # Delete a single agent
  construct agent delete coder

  # Delete multiple agents at once
  construct agent rm architect security-reviewer

  # Force delete without a confirmation prompt
  construct agent delete old-agent --force`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getAPIClient(cmd.Context())

			var agentIDs = make(map[string]string)
			for _, idOrName := range args {
				agentID, err := getAgentID(cmd.Context(), client, idOrName)
				if err != nil {
					return err
				}

				agentIDs[idOrName] = agentID
			}

			if !options.Force && !confirmDeletion(cmd.InOrStdin(), cmd.OutOrStdout(), "agent", args) {
				return nil
			}

			for idOrName, agentID := range agentIDs {
				_, err := client.Agent().DeleteAgent(cmd.Context(), &connect.Request[v1.DeleteAgentRequest]{
					Msg: &v1.DeleteAgentRequest{Id: agentID},
				})
				if err != nil {
					return fmt.Errorf("failed to delete agent %s: %w", idOrName, err)
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&options.Force, "force", "f", false, "Skip the confirmation prompt")

	return cmd
}
