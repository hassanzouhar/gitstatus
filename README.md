# Git Status Check

A powerful Go CLI tool that enhances Git workflow by providing comprehensive status checking and automated actions. This tool helps developers maintain clean Git repositories by detecting common issues and offering interactive solutions.

## Features

- 🔍 **Comprehensive Status Checks**
- Uncommitted changes detection
- Unpushed commits tracking
- Branch divergence detection
- Protected branch warnings
- Branch synchronization status

- 🤖 **Interactive Mode**
- Guided resolution of Git status issues
- Automated commit and push operations
- Interactive prompts for common Git actions

- 🎨 **User-Friendly Output**
- Color-coded status messages
- Clear action recommendations
- Progress indicators for operations

## Dependencies

- Go 1.17 or higher
- Git (installed and available in PATH)

## Installation

You can install `gitstatus` using Go's package manager:

```bash
go install github.com/your-username/git-status-check@v0.1.0
```

Or build from source:

```bash
git clone https://github.com/your-username/git-status-check.git
cd git-status-check
go build -o gitstatus
```

## Usage

### Basic Usage
```bash
gitstatus check
```

### Available Commands
```bash
gitstatus [command]

Available Commands:
check       Run a comprehensive Git status check
interactive Start interactive mode for resolving issues
help        Help about any command
```

### Flags
- `--quiet, -q`: Quiet mode (only output if there are issues)
- `--verbose, -v`: Verbose mode (show additional information)
- `--help, -h`: Help for gitstatus

### Examples

1. Check status quietly:
```bash
gitstatus check --quiet
```

2. Interactive mode with verbose output:
```bash
gitstatus interactive --verbose
```

3. Get help for a specific command:
```bash
gitstatus help check
```

## Exit Codes

The script uses different exit codes to indicate various states:

- `0`: Everything is clean (no issues)
- `1`: Script error or invalid usage
- `2`: Protected branch warning
- `3`: Uncommitted changes or unpushed commits
- `4`: Branch needs pulling (behind remote)
- `5`: Branches have diverged
- `6`: Other Git-related issues

## Color-Coded Output

The script uses color coding for better visibility:

- 🟢 Green: Success messages and completed actions
- 🔴 Red: Errors and critical issues
- 🟡 Yellow: Warnings and status notifications
- 🔵 Blue: Information and progress messages

## Interactive Mode Flow

When running in interactive mode (-i), the script:

1. Checks current branch status
2. Detects any issues (uncommitted changes, unpushed commits, etc.)
3. Offers appropriate actions for each issue
4. Guides you through resolving each issue step by step
5. Confirms successful completion of actions

Example interaction:
```
=== Recommended Actions ===
1. Commit your changes:
git add .
git commit -m "your commit message"
Would you like to commit all changes? (y/n):
```

## Contributing

Contributions are welcome! Here's how you can help:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Write and test your Go code
4. Ensure your code follows Go best practices and conventions
5. Update tests and documentation
6. Submit a Pull Request

Make sure to run tests before submitting:
```bash
go test ./...
```

## License

This project is open source and available under the MIT License.

