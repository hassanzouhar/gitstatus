# Git Status Check

A powerful bash script that enhances Git workflow by providing comprehensive status checking and automated actions. This tool helps developers maintain clean Git repositories by detecting common issues and offering interactive solutions.

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

## Installation

1. Clone this repository:
```bash
git clone https://github.com/your-username/git-status-check.git
```

2. Make the script executable:
```bash
chmod +x git_status_check.sh
```

3. (Optional) Add to your PATH for system-wide access:
```bash
ln -s $(pwd)/git_status_check.sh /usr/local/bin/git-status-check
```

## Usage

### Basic Usage
```bash
./git_status_check.sh
```

### Available Options

- `-h`: Display help information
- `-q`: Quiet mode (only output if there are issues)
- `-v`: Verbose mode (show additional information)
- `-i`: Interactive mode (guides you through resolving issues)

### Examples

1. Check status quietly:
```bash
./git_status_check.sh -q
```

2. Interactive mode with verbose output:
```bash
./git_status_check.sh -i -v
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

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open source and available under the MIT License.

