# Getting Started with Construct

Welcome to Construct! This guide will walk you through installation, setup, and your first coding session with AI agents.

## What You'll Learn

- Installing Construct and starting the daemon
- Configuring AI model providers
- Creating your first agents
- Having conversations with agents
- Basic workflows and best practices

**Time required:** 10-15 minutes

## Prerequisites

Before you begin, make sure you have:

- **Operating System:** macOS or Linux (Windows support coming soon)
- **API Key:** At least one of:
  - Anthropic API key (recommended for starting)
  - OpenAI API key
  - Google Gemini API key
  - xAI API key

Don't have an API key yet? Sign up for [Anthropic Claude](https://console.anthropic.com/) or [OpenAI](https://platform.openai.com/) to get started.

## Step 1: Installation

### Option A: Homebrew (macOS/Linux)

```bash
brew install furisto/tap/construct
```

### Option B: Download Binary

Download the latest release for your platform from the [GitHub Releases page](https://github.com/furisto/construct/releases).

```bash
# Extract and move to your PATH
tar -xzf construct_*.tar.gz
sudo mv construct /usr/local/bin/

# Verify installation
construct --help
```

### Option C: Build from Source

```bash
# Clone the repository
git clone https://github.com/furisto/construct
cd construct

# Build the CLI
cd frontend/cli
go build -o construct

# Move to your PATH (optional but recommended)
sudo mv construct /usr/local/bin/

# Verify installation
construct --help
```

## Step 2: Start the Daemon

Construct uses a background daemon to manage agents, tasks, and conversations. Install it first:

```bash
# Install the daemon as a system service
construct daemon install
```

This command:

- Creates a system service (systemd on Linux, launchd on macOS)
- Configures socket activation (daemon starts automatically when needed)
- Requires no elevated privileges (runs as your user)

**Verification:**

The daemon will start automatically when you run your first command. You can verify it's configured correctly:

```bash
# On Linux
systemctl --user status construct

# On macOS
launchctl list | grep construct
```

You should see the service listed (it may show as "not running" until first use - this is normal with socket activation).

## Step 3: Configure a Model Provider

Before creating agents, you need to connect Construct to an AI model provider. We'll start with Anthropic Claude, but you can use any supported provider.

### Using Anthropic Claude (Recommended)

**Option 1: Using environment variable**

```bash
# Set your API key as an environment variable
export ANTHROPIC_API_KEY="sk-ant-api03-..."

# Create the provider
construct modelprovider create anthropic --type anthropic

# Verify it was created
construct modelprovider list
```

**Option 2: Using the --api-key flag**

```bash
# Pass API key directly via flag
construct modelprovider create anthropic --type anthropic --api-key "sk-ant-api03-..."

# Verify
construct modelprovider list
```

**Option 3: Interactive prompt**

```bash
# Let construct prompt you for the API key
construct modelprovider create anthropic --type anthropic
# You'll be prompted: Enter Anthropic API key for anthropic:
```

You should see output like:

```
ID                                   NAME        TYPE        ENABLED
01974c1d-0be8-70e1-88b4-ad9462fff25e anthropic   anthropic   true
```

### Alternative: Using OpenAI

```bash
# Set your API key
export OPENAI_API_KEY="sk-..."

# Create the provider
construct modelprovider create openai --type openai

# Verify
construct modelprovider list
```

### Supported Providers

Construct supports these providers (you can add multiple):

- **anthropic** - Claude models (Opus, Sonnet, Haiku)
- **openai** - GPT models (GPT-4, GPT-3.5)
- **gemini** - Google Gemini models
- **xai** - Grok models

**Tip:** You can configure multiple providers and switch between them as needed.

## Step 4: Register Models

Now register the specific models you want to use. Each model needs a name, provider, and context window size.

### For Anthropic Users

```bash
# Register Claude 3.5 Sonnet (balanced - recommended for starting)
construct model create claude-3-5-sonnet-20241022 \
  --provider anthropic \
  --context-window 200000

# Register Claude 3.5 Haiku (fast - for quick tasks)
construct model create claude-3-5-haiku-20241022 \
  --provider anthropic \
  --context-window 200000

# Verify your models
construct model list
```

### For OpenAI Users

```bash
# Register GPT-4o (capable)
construct model create gpt-4o \
  --provider openai \
  --context-window 128000

# Register GPT-4o mini (fast)
construct model create gpt-4o-mini \
  --provider openai \
  --context-window 128000

# Verify
construct model list
```

You should see your models listed:

```
ID                                   NAME                        PROVIDER    CONTEXT_WINDOW  ENABLED
01974c1d-0be8-70e1-88b4-ad9462fff25e claude-3-5-sonnet-20241022  anthropic   200000          true
01974c1d-0be8-70e1-88b4-ad9462fff26f claude-3-5-haiku-20241022   anthropic   200000          true
```

## Step 5: Create Your First Agent

Agents are AI assistants with specific instructions and capabilities. Let's create a general-purpose coding agent:

```bash
construct agent create edit \
  --prompt "You are a helpful coding assistant who writes clean, well-documented code and explains your reasoning." \
  --model claude-3-5-sonnet-20241022 \
  --description "General purpose coding assistant"
```

**Success!** You should see:

```
Created agent: edit (01974c1d-0be8-70e1-88b4-ad9462fff25e)
```

### Verify Your Agent

```bash
construct agent list
```

You should see:

```
ID                                   NAME  MODEL                       ENABLED
01974c1d-0be8-70e1-88b4-ad9462fff25e edit  claude-3-5-sonnet-20241022  true
```

### Optional: Create a Quick Agent

For rapid iterations, create a fast agent using a lighter model:

```bash
construct agent create quick \
  --prompt "You are a fast coding assistant for quick tasks like writing commit messages, finding files, or making small changes." \
  --model claude-3-5-haiku-20241022 \
  --description "Fast agent for simple tasks"
```

## Step 6: Have Your First Conversation

Now for the fun part - let's talk to your agent!

### Interactive Mode

Start an interactive conversation:

```bash
construct new --agent edit
```

You'll see a prompt where you can chat with your agent:

```
Starting new conversation with agent: edit
Type your message (Ctrl+C to exit):

> 
```

Try asking something like:

```
> Write a Python function that calculates the factorial of a number
```

The agent will respond with code and explanation. You can continue the conversation:

```
> Now add error handling for negative numbers
```

**Exit:** Press Ctrl+C when you're done. Your conversation is automatically saved!

### Quick Questions

For one-off questions, use the ask command:

```bash
construct ask "What is the time complexity of bubble sort?"
```

The agent responds and exits. Perfect for quick queries!

### With File Context

Ask questions about your code:

```bash
# Ask about a specific file
construct ask "What does this code do?" --file main.go

# Include multiple files
construct ask "Are these interfaces consistent?" \
  --file user.go \
  --file product.go

# Pipe input
cat error.log | construct ask "What's causing this error?"
```

## Step 7: Resume Your Work

Every conversation is saved as a task. Resume where you left off:

```bash
# Show recent tasks and pick one
construct resume

# Or resume the most recent task immediately
construct resume --last
```

Your full conversation history is restored, including:

- All previous messages
- File context
- The agent you were using
- Workspace directory

**Pro tip:** You can even switch agents when resuming to get a different perspective!

```bash
# Start with edit agent
construct new --agent edit
# ... work for a while ...

# Resume with a different agent
construct resume --last --agent quick
```

## Next Steps

Congratulations! You've completed the basic setup. Here's what to explore next:

- **[CLI Reference](cli_reference.md)** - Complete command reference and advanced usage
- **[Architecture](architecture.md)** - Understanding how Construct works under the hood
- **[Troubleshooting](troubleshooting.md)** - Common issues and solutions

## Example: Complete Setup Script

Here's a complete setup script you can run:

```bash
#!/bin/bash

# Install daemon
construct daemon install

# Set up Anthropic (alternatively, use --api-key flag)
export ANTHROPIC_API_KEY="sk-ant-..."
construct modelprovider create anthropic --type anthropic

# Register models
construct model create claude-3-5-sonnet-20241022 \
  --provider anthropic \
  --context-window 200000

construct model create claude-3-5-haiku-20241022 \
  --provider anthropic \
  --context-window 200000

# Create agents
construct agent create edit \
  --prompt "You are a helpful coding assistant who writes clean, well-documented code." \
  --model claude-3-5-sonnet-20241022

construct agent create quick \
  --prompt "You are a fast coding assistant for quick tasks." \
  --model claude-3-5-haiku-20241022

# Set defaults (optional)
construct config set cmd.new.agent edit
construct config set cmd.ask.agent quick

echo "Setup complete! Try: construct new --agent edit"
```

---

You're all set! Start building with Construct:

```bash
construct new --agent edit
```

Welcome to a new way of coding with AI.
