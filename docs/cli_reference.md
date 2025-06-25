# Construct CLI Reference

The Construct CLI is a command-line interface for managing AI agents, tasks, messages, models, and model providers. This document provides a complete reference for all available commands and their options.

## Global Options

- `-v, --verbose`: Enable verbose output
- `--help`: Show help information for any command

## Commands Overview

The CLI commands are organized into three main groups:

- **Core Commands**: Primary workflow commands for interactive usage
- **Resource Management**: Commands for managing agents, tasks, messages, models, and providers
- **System Commands**: Configuration, daemon management, and utilities

---

## Core Commands

### `construct new`

Start a new interactive conversation with an AI agent.

**Usage:**
```bash
construct new [flags]
```

**Flags:**
- `--agent string`: Use a specific agent (default: last used or configured default)
- `--workspace string`: The sandbox directory for the agent (default: current directory)

**Examples:**
```bash
# Start a new conversation with the default agent
construct new

# Start with a specific agent
construct new --agent coder

# Sandbox another directory
construct new --workspace /workspace/repo/hello/world
```

### `construct resume`

Resume an existing conversation from where you left off.

**Usage:**
```bash
construct resume [task-id] [flags]
```

**Description:**
Resume an existing conversation, restoring full context including all previous messages, the agent that was being used, and workspace directory settings. If no task ID is provided, shows an interactive picker of recent tasks. Supports partial ID matching for convenience.

**Flags:**
- `--last`: Resume the most recent task immediately

**Examples:**
```bash
# Show interactive picker to select from recent tasks
construct resume

# Resume the most recent task immediately
construct resume --last

# Resume specific task by full ID
construct resume 01974c1d-0be8-70e1-88b4-ad9462fff25e
```

### `construct ask`

Ask a question to the AI and create a saved task that can be resumed later.

**Usage:**
```bash
construct ask [question] [flags]
```

**Flags:**
- `-a, --agent string`: The agent to use (name or ID)
- `-w, --workspace string`: The workspace directory
- `--max-turns int`: Maximum number of turns for the conversation (default: 5)
- `-f, --file strings`: Files to include as context (can be used multiple times)
- `-c, --continue`: Continue the previous task

**Examples:**
```bash
# Simple question
construct ask "What is 2+2?"

# Use a specific agent
construct ask "Review this code for security issues" --agent security-reviewer

# Include files as context
construct ask "What does this code do?" --file main.go --file utils.go

# Pipe input with question and file context
cat main.go | construct ask "What does this code do?" --file config.yaml

# Give agent more turns for complex tasks
construct ask "Debug why the tests are failing" --max-turns 10 --file test.log

# Get JSON output for scripting
construct ask "List all Go files" --output json

# Complex analysis with file context
construct ask "Analyze the architecture and suggest improvements" --max-turns 15 --agent architect --file architecture.md
```

---

## Resource Management Commands

### Agent Commands

#### `construct agent create`

Create a new AI agent with custom instructions and model configuration.

**Usage:**
```bash
construct agent create <name> [flags]
```

**Arguments:**
- `name` (required): Name of the agent

**Flags:**
- `-d, --description string`: Description of the agent
- `-p, --prompt string`: System prompt that defines the agent's behavior
- `--prompt-file string`: Read system prompt from file
- `--prompt-stdin`: Read system prompt from stdin
- `-m, --model string`: AI model to use (e.g. gpt-4o, claude-4 or model ID) (required)

**Note:** Exactly one prompt source must be specified (--prompt, --prompt-file, or --prompt-stdin).

**Examples:**
```bash
construct agent create "coder" --prompt "You are a coding assistant" --model "claude-4"
construct agent create "sql-expert" --prompt-file ./prompts/sql-expert.txt --model "claude-4"
echo "You review code" | construct agent create "reviewer" --prompt-stdin --model "gpt-4o"
construct agent create "RFC writer" --prompt "You help with writing" --model "gemini-2.5.pro" --description "RFC writing assistant"
```

#### `construct agent list`

List all available agents with optional filtering.

**Usage:**
```bash
construct agent list [flags]
```

**Aliases:** `construct agent ls`

**Flags:**
- `-m, --model string`: Show only agents using this AI model (e.g., 'claude-4', 'gpt-4', or model ID)
- `-n, --name string`: Filter agents by name
- `-l, --limit int`: Limit number of results
- `--enabled`: Show only enabled agents (default: true)
- `--output string`: Output format (json, yaml, table)

**Examples:**
```bash
construct agent list
construct agent list --model "claude-4"
construct agent list --enabled
construct agent list --model "claude-4" --enabled --limit 5
```

#### `construct agent get`

Get detailed information about a specific agent.

**Usage:**
```bash
construct agent get <id-or-name> [flags]
```

**Arguments:**
- `id-or-name` (required): Agent ID or name

**Flags:**
- `--output string`: Output format (json, yaml, table)

**Examples:**
```bash
construct agent get "coder"
construct agent get 01974c1d-0be8-70e1-88b4-ad9462fff25e
construct agent get "sql-expert" --output json
construct agent get "reviewer" --output yaml
```

#### `construct agent edit`

Edit an agent configuration interactively using your default editor.

**Usage:**
```bash
construct agent edit <id-or-name>
```

**Description:**
Edit an agent configuration using your default editor ($EDITOR). This command fetches the current agent configuration, opens it as a YAML file in your terminal editor, and applies any changes you make upon saving and closing the editor. Similar to 'kubectl edit', this provides a powerful way to make multiple changes to an agent in a single operation.

**Arguments:**
- `id-or-name` (required): Agent ID or name

**Examples:**
```bash
construct agent edit "coder"
construct agent edit 01974c1d-0be8-70e1-88b4-ad9462fff25e
EDITOR=nano construct agent edit "coder"
```

#### `construct agent delete`

Delete one or more agents by their IDs or names.

**Usage:**
```bash
construct agent delete <id-or-name>... [flags]
```

**Aliases:** `construct agent rm`

**Arguments:**
- `id-or-name...` (required): One or more agent IDs or names

**Flags:**
- `-f, --force`: Force deletion without confirmation

**Examples:**
```bash
construct agent delete coder architect debugger
construct agent delete 01974c1d-0be8-70e1-88b4-ad9462fff25e
construct agent delete coder --force
```

### Task Commands

#### `construct task create`

Create a new task and assign it to an agent.

**Usage:**
```bash
construct task create [flags]
```

**Flags:**
- `-a, --agent string`: The agent to assign to the task (name or ID) (required)
- `-w, --workspace string`: The workspace directory

**Examples:**
```bash
construct task create --agent coder
construct task create --agent sql-expert --workspace /path/to/repo
```

#### `construct task list`

List all tasks with optional filtering.

**Usage:**
```bash
construct task list [flags]
```

**Aliases:** `construct task ls`

**Flags:**
- `-a, --agent string`: Filter by agent (name or ID)
- `--output string`: Output format (json, yaml, table)

**Examples:**
```bash
construct task list
construct task list --agent "coder"
construct task list --output yaml
```

#### `construct task get`

Get detailed information about a specific task.

**Usage:**
```bash
construct task get <task-id> [flags]
```

**Arguments:**
- `task-id` (required): Task ID

**Flags:**
- `--output string`: Output format (json, yaml, table)

**Examples:**
```bash
construct task get 01974c1d-0be8-70e1-88b4-ad9462fff25e
construct task get 01974c1d-0be8-70e1-88b4-ad9462fff25e --output json
```

#### `construct task delete`

Delete one or more tasks by their IDs.

**Usage:**
```bash
construct task delete <task-id>... [flags]
```

**Aliases:** `construct task rm`

**Arguments:**
- `task-id...` (required): One or more task IDs

**Flags:**
- `-f, --force`: Force deletion without confirmation

**Examples:**
```bash
construct task delete 01974c1d-0be8-70e1-88b4-ad9462fff25e
construct task delete 01974c1d-0be8-70e1-88b4-ad9462fff25e 01974c1d-0be8-70e1-88b4-ad9462fff26f
```

### Message Commands

#### `construct message create`

Create a new message for a task.

**Usage:**
```bash
construct message create <task-id> <content> [flags]
```

**Arguments:**
- `task-id` (required): Task ID
- `content` (required): Message content

**Flags:**
- `--output string`: Output format (json, yaml, table)

**Examples:**
```bash
construct message create "123e4567-e89b-12d3-a456-426614174000" "Please implement a hello world function"
```

#### `construct message list`

List messages with optional filtering.

**Usage:**
```bash
construct message list [flags]
```

**Aliases:** `construct message ls`, `construct msg ls`

**Flags:**
- `-t, --task string`: Filter by task ID
- `-a, --agent string`: Filter by agent name or ID
- `-r, --role string`: Filter by role (user or assistant)
- `--output string`: Output format (json, yaml, table)

**Examples:**
```bash
construct message list
construct message list --agent "coder"
construct message list --task "456e7890-e12b-34c5-a678-901234567890"
construct message list --role assistant
```

#### `construct message get`

Get detailed information about a specific message.

**Usage:**
```bash
construct message get <message-id> [flags]
```

**Arguments:**
- `message-id` (required): Message ID

**Flags:**
- `--output string`: Output format (json, yaml, table)

**Examples:**
```bash
construct message get "123e4567-e89b-12d3-a456-426614174000"
```

#### `construct message delete`

Delete one or more messages by ID.

**Usage:**
```bash
construct message delete <message-id>... [flags]
```

**Aliases:** `construct message rm`, `construct msg rm`

**Arguments:**
- `message-id...` (required): One or more message IDs

**Flags:**
- `-f, --force`: Force deletion without confirmation

**Examples:**
```bash
construct message delete "123e4567-e89b-12d3-a456-426614174000"
```

### Model Commands

#### `construct model create`

Create a new model with specified configuration.

**Usage:**
```bash
construct model create <model-name> [flags]
```

**Arguments:**
- `model-name` (required): Name of the model

**Flags:**
- `-p, --provider string`: The name or ID of the model provider (required)
- `-w, --context-window int`: The context window size (required)

**Examples:**
```bash
construct model create "gpt-4" --provider "openai-dev" --context-window 8192
construct model create "claude-3-5-sonnet" --provider "123e4567-e89b-12d3-a456-426614174000" --context-window 200000
```

#### `construct model list`

List all available models with optional filtering.

**Usage:**
```bash
construct model list [flags]
```

**Aliases:** `construct model ls`

**Flags:**
- `-p, --provider string`: Filter by model provider name or ID
- `-d, --show-disabled`: Show disabled models
- `--output string`: Output format (json, yaml, table)

**Examples:**
```bash
construct model list
construct model list --provider "anthropic-dev"
construct model list --show-disabled
```

#### `construct model get`

Get detailed information about a specific model.

**Usage:**
```bash
construct model get <model-id-or-name> [flags]
```

**Arguments:**
- `model-id-or-name` (required): Model ID or name

**Flags:**
- `--output string`: Output format (json, yaml, table)

**Examples:**
```bash
construct model get "gpt-4"
construct model get "123e4567-e89b-12d3-a456-426614174000"
```

#### `construct model delete`

Delete one or more models by ID or name.

**Usage:**
```bash
construct model delete <model-id-or-name>... [flags]
```

**Aliases:** `construct model rm`

**Arguments:**
- `model-id-or-name...` (required): One or more model IDs or names

**Flags:**
- `-f, --force`: Force deletion without confirmation

**Examples:**
```bash
construct model delete "gpt-4"
construct model delete "claude-3-5-sonnet" "llama-3.1-8b" "gpt-4"
```

### Model Provider Commands

#### `construct modelprovider create`

Create a new model provider integration.

**Usage:**
```bash
construct modelprovider create <name> [flags]
```

**Description:**
Configure integrations to AI model providers to access their language models for your agents. Providers require API credentials and offer different model capabilities. At least one provider must be configured before creating agents.

**Supported providers:**
- **OpenAI**: Access to GPT models (gpt-4, gpt-3.5-turbo, etc.)
- **Anthropic**: Access to Claude models (claude-3-5-sonnet, claude-3-haiku, etc.)

**Arguments:**
- `name` (required): Name of the provider

**Flags:**
- `-k, --api-key string`: The API key for the model provider (can also be set via environment variable)
- `-t, --type string`: The type of the model provider (anthropic, openai) (required)

**Examples:**
```bash
construct modelprovider create "openai-dev" --type openai
export OPENAI_API_KEY="sk-..." && construct modelprovider create "openai-prod" --type openai
construct modelprovider create "anthropic-prod" --type anthropic --api-key "sk-ant-..."
```

#### `construct modelprovider list`

List all model providers with optional filtering.

**Usage:**
```bash
construct modelprovider list [flags]
```

**Aliases:** `construct modelprovider ls`, `construct mp ls`

**Flags:**
- `-t, --provider-type strings`: Filter by provider type (anthropic, openai)
- `--enabled`: Show only enabled model providers (default: true)
- `--output string`: Output format (json, yaml, table)

**Examples:**
```bash
construct modelprovider list --enabled=false
construct modelprovider list --provider-type anthropic --provider-type openai
```

#### `construct modelprovider get`

Get detailed information about a specific model provider.

**Usage:**
```bash
construct modelprovider get <id-or-name> [flags]
```

**Arguments:**
- `id-or-name` (required): Provider ID or name

**Flags:**
- `--output string`: Output format (json, yaml, table)

**Examples:**
```bash
construct modelprovider get "anthropic-dev"
construct modelprovider get "openai-prod" --output json
```

#### `construct modelprovider delete`

Delete one or more model providers by their IDs or names.

**Usage:**
```bash
construct modelprovider delete <id-or-name>... [flags]
```

**Aliases:** `construct modelprovider rm`, `construct mp rm`

**Arguments:**
- `id-or-name...` (required): One or more provider IDs or names

**Flags:**
- `-f, --force`: Force deletion without confirmation

**Note:** This operation also deletes all associated models.

**Examples:**
```bash
construct modelprovider delete anthropic-dev openai-prod
construct modelprovider delete 01974c1d-0be8-70e1-88b4-ad9462fff25e
```

---

## System Commands

### Configuration Commands

#### `construct config set`

Set a configuration value using dot notation for nested keys.

**Usage:**
```bash
construct config set <key> <value>
```

**Arguments:**
- `key` (required): Configuration key (supports dot notation)
- `value` (required): Configuration value

**Supported Configuration Keys:**
- `cmd.new.agent`: Default agent for new command
- `cmd.ask.agent`: Default agent for ask command  
- `cmd.ask.max-turns`: Default max turns for ask command
- `cmd.resume.recent_task_limit`: Recent task limit for resume command
- `log.level`: Logging level
- `log.file`: Log file path
- `log.format`: Log format
- `editor`: Default editor
- `output.format`: Default output format
- `output.no-headers`: Disable headers in output
- `output.wide`: Enable wide output format

**Examples:**
```bash
construct config set cmd.new.agent "coder"
construct config set output.format "json"
construct config set log.level "debug"
```

#### `construct config get`

Get a configuration value.

**Usage:**
```bash
construct config get <key>
```

**Arguments:**
- `key` (required): Configuration key

**Examples:**
```bash
construct config get cmd.new.agent
construct config get output.format
```

#### `construct config list`

List all configuration values.

**Usage:**
```bash
construct config list
```

#### `construct config unset`

Remove a configuration value.

**Usage:**
```bash
construct config unset <key> [flags]
```

**Arguments:**
- `key` (required): Configuration key

**Flags:**
- `-f, --force`: Force the removal without confirmation

**Examples:**
```bash
construct config unset cmd.new.agent
construct config unset output.format --force
```

### Daemon Commands

#### `construct daemon install`

Install the Construct daemon as a system service.

**Usage:**
```bash
construct daemon install [flags]
```

**Description:**
Install the daemon with platform-specific service management (launchd on macOS, systemd on Linux).

**Flags:**
- `-f, --force`: Force install the daemon
- `--always-running`: Run the daemon continuously instead of using socket activation
- `--listen-http string`: HTTP address to listen on
- `-q, --quiet`: Silent installation
- `-n, --name string`: Name of the daemon (used for socket activation and context) (default: "default")

**Examples:**
```bash
construct daemon install
construct daemon install --listen-http 127.0.0.1:8080
construct daemon install --force
construct daemon install --name production --listen-http :8080 --always-running
```

#### `construct daemon run`

Run the API server as a persistent service.

**Usage:**
```bash
construct daemon run [flags]
```

**Description:**
Run the construct server as a single, long-running process. Supports different launch modes:

**On macOS:**
- If launched by launchd: uses HTTP address if provided, otherwise uses socket activation
- If not launched by launchd: uses provided HTTP address or Unix socket

**On Linux:**
- If launched by systemd: uses HTTP address if provided, otherwise uses socket activation  
- If not launched by systemd: uses provided HTTP address or Unix socket

**Flags:**
- `--listen-http string`: The address to listen on for HTTP requests
- `--listen-unix string`: The path to listen on for Unix socket requests

#### `construct daemon stop`

Stop the running daemon.

**Usage:**
```bash
construct daemon stop
```

#### `construct daemon uninstall`

Uninstall the Construct daemon from the system.

**Usage:**
```bash
construct daemon uninstall [flags]
```

**Flags:**
- `-y, --yes`: Skip confirmation prompt
- `-q, --quiet`: Quiet mode

**Examples:**
```bash
construct daemon uninstall
construct daemon uninstall -y
```

### Utility Commands

#### `construct version`

Print the version number of Construct.

**Usage:**
```bash
construct version
```

#### `construct update`

Update the CLI to the latest version.

**Usage:**
```bash
construct update
```

---

## Common Patterns

### Output Formats

Many commands support multiple output formats:
- `table` (default): Human-readable table format
- `json`: JSON format for scripting and automation
- `yaml`: YAML format for configuration files

Use the `--output` flag to specify the format:
```bash
construct agent list --output json
construct task get <task-id> --output yaml
```

### ID and Name Resolution

Commands that accept IDs or names automatically resolve names to IDs when needed. You can use either:
- Full UUIDs: `01974c1d-0be8-70e1-88b4-ad9462fff25e`
- Human-readable names: `coder`, `sql-expert`

### Confirmation Prompts

Delete operations include confirmation prompts unless the `--force` flag is used:
```bash
construct agent delete coder        # Shows confirmation prompt
construct agent delete coder --force # Skips confirmation
```

### Aliases

Many commands have shorter aliases for convenience:
- `construct agent ls` → `construct agent list`
- `construct task rm` → `construct task delete`
- `construct modelprovider` → `construct mp`
- `construct message` → `construct msg`

### Context and Configuration

The CLI automatically manages contexts and configuration:
- Configuration is stored in `~/.construct/config.yaml`
- Context information is stored in `~/.construct/context.yaml`
- The daemon must be installed before using resource management commands

---

## Getting Help

For detailed help on any command, use the `--help` flag:
```bash
construct --help                    # General help
construct agent --help              # Help for agent commands
construct agent create --help       # Help for specific subcommand
```

For more information and documentation, visit the project repository or documentation site.