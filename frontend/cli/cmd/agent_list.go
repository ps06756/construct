package cmd

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	api "github.com/furisto/construct/api/go/client"
	v1 "github.com/furisto/construct/api/go/v1"
	"github.com/furisto/construct/frontend/cli/pkg/fail"
	"github.com/spf13/cobra"
)

type agentListOptions struct {
	Models        []string
	Names         []string
	Limit         int32
	RenderOptions RenderOptions
}

func NewAgentListCmd() *cobra.Command {
	var options agentListOptions

	cmd := &cobra.Command{
		Use:     "list [flags]",
		Short:   "List all available agents",
		Aliases: []string{"ls"},
		Example: `  # List all agents in a table
  construct agent list

  # Find all agents using a specific model
  construct agent ls --model "claude-3-5-sonnet"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getAPIClient(cmd.Context())

			agents, err := agentList(cmd.Context(), options, client)
			if err != nil {
				return fail.HandleError(err)
			}

			return getRenderer(cmd.Context()).Render(agents, &options.RenderOptions)
		},
	}

	cmd.Flags().StringArrayVarP(&options.Models, "model", "m", []string{}, "Filter agents by the model they use")
	cmd.Flags().StringArrayVarP(&options.Names, "name", "n", []string{}, "Filter agents by name (supports partial matching)")
	cmd.Flags().Int32VarP(&options.Limit, "limit", "l", 0, "Limit the number of results returned")
	addRenderOptions(cmd, &options.RenderOptions)

	return cmd
}

func agentList(ctx context.Context, options agentListOptions, client *api.Client) ([]*AgentDisplay, error) {
	filter := &v1.ListAgentsRequest_Filter{}

	if len(options.Names) > 0 {
		filter.Names = options.Names
	}

	// Handle enabled/disabled filtering
	// Note: API doesn't currently support enabled/disabled filtering
	// This would need to be implemented on the backend first

	req := &connect.Request[v1.ListAgentsRequest]{
		Msg: &v1.ListAgentsRequest{
			Filter: filter,
		},
	}

	// Note: API doesn't currently support limit parameter
	// This would need to be implemented on the backend first

	resp, err := client.Agent().ListAgents(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list agents: %w", err)
	}

	displayAgents := make([]*AgentDisplay, len(resp.Msg.Agents))
	for i, agent := range resp.Msg.Agents {
		model, err := client.Model().GetModel(ctx, &connect.Request[v1.GetModelRequest]{
			Msg: &v1.GetModelRequest{
				Id: agent.Spec.ModelId,
			},
		})

		if err != nil {
			return nil, fmt.Errorf("failed to get model %s: %w", agent.Spec.ModelId, err)
		}

		displayAgents[i] = ConvertAgentToDisplay(agent, model.Msg.Model.Spec.Name)
	}

	return displayAgents, nil
}
