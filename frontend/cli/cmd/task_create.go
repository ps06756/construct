package cmd

import (
	"fmt"

	"connectrpc.com/connect"
	v1 "github.com/furisto/construct/api/go/v1"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type taskCreateOptions struct {
	Agent            string
	ProjectDirectory string
}

func NewTaskCreateCmd() *cobra.Command {
	var options taskCreateOptions

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new task",
		Example: `  # Create a basic task
  construct task create

  # Create task and assign to agent by name
  construct task create --agent "coder"

  # Create task and assign to agent by ID
  construct task create --agent 01974c1d-0be8-70e1-88b4-ad9462fff25e

  # Create task with specific project directory
  construct task create --project-directory /path/to/project

  # Create task with both agent and project directory
  construct task create --agent "sql-expert" --project-directory /path/to/project

  # Short form flags
  construct task create -a "reviewer" -d /path/to/project`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getAPIClient(cmd.Context())

			agentID := options.Agent
			if agentID != "" {
				_, err := uuid.Parse(agentID)
				if err != nil {
					resolvedID, err := getAgentID(cmd.Context(), client, agentID)
					if err != nil {
						return fmt.Errorf("failed to resolve agent %s: %w", agentID, err)
					}
					agentID = resolvedID
				}
			}

			resp, err := client.Task().CreateTask(cmd.Context(), &connect.Request[v1.CreateTaskRequest]{
				Msg: &v1.CreateTaskRequest{
					AgentId:          agentID,
					ProjectDirectory: options.ProjectDirectory,
				},
			})

			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), resp.Msg.Task.Id)
			return nil
		},
	}

	cmd.Flags().StringVarP(&options.Agent, "agent", "a", "", "The agent to assign to the task (name or ID)")
	cmd.Flags().StringVarP(&options.ProjectDirectory, "project-directory", "d", "", "The project directory")

	return cmd
}
