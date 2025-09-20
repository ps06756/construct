package cmd

import (
	"fmt"
	"io"

	"connectrpc.com/connect"

	v1 "github.com/furisto/construct/api/go/v1"
	"github.com/google/uuid"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type agentCreateOptions struct {
	Description  string
	SystemPrompt string
	PromptFile   string
	PromptStdin  bool
	Model        string
}

func NewAgentCreateCmd() *cobra.Command {
	var options agentCreateOptions

	var cmd = &cobra.Command{
		Use:   "create <name> [flags]",
		Short: "Define a new, reusable AI agent",
		Args:  cobra.ExactArgs(1),
		Long: `Define a new, reusable AI agent.

Creates a new agent by giving it a name, a system prompt, and assigning a model. 
The system prompt defines the agent's personality, goals, and constraints.

You must specify the agent's system prompt using one of --prompt, --prompt-file, 
or by piping it via --prompt-stdin.`,
		Example: `  # Create a simple coding assistant
  construct agent create "coder" \
    --model "gpt-4o" \
    --prompt "You are an expert Go developer. Your code is clean, efficient, and well-documented."

  # Create an agent with a prompt from a file
  construct agent create "sql-expert" \
    --model "claude-3-5-sonnet" \
    --prompt-file ./prompts/sql.txt

  # Create an agent by piping the prompt
  echo "You are a security expert reviewing code for vulnerabilities." | \
    construct agent create "reviewer" --model "gpt-4o" --prompt-stdin`,
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			systemPrompt, err := getSystemPrompt(&options, cmd.InOrStdin(), getFileSystem(cmd.Context()))
			if err != nil {
				return err
			}

			client := getAPIClient(cmd.Context())

			_, err = uuid.Parse(options.Model)
			if err != nil {
				modelID, err := getModelID(cmd.Context(), client, options.Model)
				if err != nil {
					return err
				}
				options.Model = modelID
			}

			agentResp, err := client.Agent().CreateAgent(cmd.Context(), &connect.Request[v1.CreateAgentRequest]{
				Msg: &v1.CreateAgentRequest{
					Name:         name,
					Description:  options.Description,
					Instructions: systemPrompt,
					ModelId:      options.Model,
				},
			})

			if err != nil {
				return fmt.Errorf("failed to create agent: %w", err)
			}

			cmd.Println(agentResp.Msg.Agent.Metadata.Id)
			return nil
		},
	}

	cmd.Flags().StringVarP(&options.Description, "description", "d", "", "A brief description of what the agent does")
	cmd.Flags().StringVarP(&options.SystemPrompt, "prompt", "p", "", "The system prompt that defines the agent's behavior")
	cmd.Flags().StringVar(&options.PromptFile, "prompt-file", "", "Read the system prompt from a specified file")
	cmd.Flags().BoolVar(&options.PromptStdin, "prompt-stdin", false, "Read the system prompt from standard input (stdin)")
	cmd.Flags().StringVarP(&options.Model, "model", "m", "", "The AI model the agent will use (e.g., gpt-4o) (required)")

	cmd.MarkFlagRequired("model")

	return cmd
}

func getSystemPrompt(options *agentCreateOptions, stdin io.Reader, fs *afero.Afero) (string, error) {
	promptSources := 0

	if options.SystemPrompt != "" {
		promptSources++
	}
	if options.PromptFile != "" {
		promptSources++
	}
	if options.PromptStdin {
		promptSources++
	}

	if promptSources == 0 {
		return "", fmt.Errorf("system prompt is required (use --prompt, --prompt-file, or --prompt-stdin)")
	}
	if promptSources > 1 {
		return "", fmt.Errorf("only one prompt source can be specified (--prompt, --prompt-file, or --prompt-stdin)")
	}

	// Inline prompt
	if options.SystemPrompt != "" {
		return options.SystemPrompt, nil
	}

	// From file
	if options.PromptFile != "" {
		content, err := fs.ReadFile(options.PromptFile)
		if err != nil {
			return "", fmt.Errorf("failed to read prompt file %s: %w", options.PromptFile, err)
		}
		if len(content) == 0 {
			return "", fmt.Errorf("prompt file %s is empty", options.PromptFile)
		}
		return string(content), nil
	}

	// From stdin
	if options.PromptStdin {
		content, err := io.ReadAll(stdin)
		if err != nil {
			return "", fmt.Errorf("failed to read prompt from stdin: %w", err)
		}
		if len(content) == 0 {
			return "", fmt.Errorf("no prompt content received from stdin")
		}
		return string(content), nil
	}

	return "", fmt.Errorf("no prompt source specified")
}
