# Claude CLI Provider

This provider allows Crush to use the local Claude CLI executable, enabling use of your Claude subscription (Basic, Plus, or Max) instead of paying for API access.

## Prerequisites

- Claude CLI must be installed and available in your PATH
- You must be logged into the Claude CLI with your subscription

## Configuration

Add the following to your Crush configuration file (`~/.config/crush/config.json` or `crush.json` in your project):

```json
{
  "providers": {
    "claude-cli": {
      "id": "claude-cli",
      "name": "Claude Subscription",
      "type": "claude-cli",
      "models": [
        {
          "id": "opus",
          "name": "Opus 4.5 (Subscription)",
          "context_window": 200000,
          "default_max_tokens": 8192,
          "supports_images": true
        },
        {
          "id": "sonnet",
          "name": "Sonnet 4 (Subscription)",
          "context_window": 200000,
          "default_max_tokens": 8192,
          "supports_images": true
        },
        {
          "id": "haiku",
          "name": "Haiku 3.5 (Subscription)",
          "context_window": 200000,
          "default_max_tokens": 8192,
          "supports_images": true
        }
      ]
    }
  },
  "models": {
    "large": {
      "model": "opus",
      "provider": "claude-cli"
    },
    "small": {
      "model": "haiku",
      "provider": "claude-cli"
    }
  }
}
```

## Provider Options

### Custom Executable Path

If your Claude CLI is not in your PATH, you can specify the path using `base_url`:

```json
{
  "providers": {
    "claude-cli": {
      "type": "claude-cli",
      "base_url": "/path/to/claude",
      ...
    }
  }
}
```

## How It Works

The Claude CLI provider:

1. Spawns the `claude` executable with `--print --output-format stream-json` flags
2. Sends the prompt via stdin
3. Parses the streaming JSON output and converts it to the fantasy format
4. Uses `--dangerously-skip-permissions` since Crush has its own permission system

## Limitations

- Tool execution is handled by Crush, not the Claude CLI's built-in tool system
- Session continuity is not currently supported (each request starts fresh)
- The provider relies on the Claude CLI's output format, which may change in future versions

## Troubleshooting

### Claude CLI not found

Make sure the `claude` executable is in your PATH:

```bash
which claude
```

If not, install Claude CLI or specify the full path in the config.

### Authentication errors

Make sure you're logged into the Claude CLI:

```bash
claude login
```

### Permission errors

The provider uses `--dangerously-skip-permissions` by default. If you need to use Claude CLI's permission system instead of Crush's, this would require modifications to the provider.
