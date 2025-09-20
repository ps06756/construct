package cmd

import (
	"fmt"
	"path/filepath"

	"connectrpc.com/connect"
	v1 "github.com/furisto/construct/api/go/v1"
	"github.com/google/uuid"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type taskCreateOptions struct {
	Agent     string
	Workspace string
}

func NewTaskCreateCmd() *cobra.Command {
	var options taskCreateOptions

	cmd := &cobra.Command{
		Use:   "create [flags]",
		Short: "Create a new task without starting an interactive session",
		Long: `Create a new task without starting an interactive session.

Programmatically creates a new task. This is useful for setting up tasks that 
will be used later or by automated systems. To start an interactive task, 
use construct new.`,
		Example: `  # Create a new task assigned to the 'coder' agent
  construct task create --agent coder

  # Create a task with a specific workspace
  construct task create --agent sql-expert --workspace /path/to/db/repo`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := getAPIClient(cmd.Context())
			fs := getFileSystem(cmd.Context())

			if options.Workspace != "" {
				exists, err := afero.Exists(fs, options.Workspace)
				if err != nil {
					return fmt.Errorf("failed to check if workspace directory %s exists: %w", options.Workspace, err)
				}
				if !exists {
					return fmt.Errorf("workspace directory %s does not exist", options.Workspace)
				}

				absPath, err := filepath.Abs(options.Workspace)
				if err != nil {
					return fmt.Errorf("failed to get absolute path of workspace directory %s: %w", options.Workspace, err)
				}

				options.Workspace = absPath
			}

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

			req := &connect.Request[v1.CreateTaskRequest]{
				Msg: &v1.CreateTaskRequest{
					AgentId:          agentID,
					ProjectDirectory: options.Workspace,
				},
			}

			resp, err := client.Task().CreateTask(cmd.Context(), req)
			if err != nil {
				return fmt.Errorf("failed to create task: %w", err)
			}

			cmd.Println(resp.Msg.Task.Metadata.Id)
			return nil
		},
	}

	cmd.Flags().StringVarP(&options.Agent, "agent", "a", "", "The agent to assign to the task (required)")
	cmd.Flags().StringVarP(&options.Workspace, "workspace", "w", "", "The workspace directory for the task")

	cmd.MarkFlagRequired("agent")

	return cmd
}
