# Gitstatus

A command-line tool that helps maintain repository hygiene by checking Git repository status and providing interactive fixes for common issues.

## Features

- 🔍 Comprehensive repository status checks:
- Uncommitted changes detection
- Unpushed commits detection
- Branch protection checks
- Update requirements check
- 🤖 Interactive mode with guided fixes
- 🚦 Clear status indicators and colored output
- 🔧 Automatic issue resolution (with confirmation)
- 🤫 Quiet mode for CI/CD pipelines

## Installation

### Prerequisites
- Go 1.18 or later
- Git

### Install from source
```bash
go install github.com/hassanzouhar/gitstatus@latest
```

Make sure `$HOME/go/bin` is in your PATH:
```bash
export PATH=$PATH:$HOME/go/bin
```

## Usage

### Basic Usage
```bash
gitstatus
```

### Interactive Mode
```bash
gitstatus -i
```

### Flags
- `-i, --interactive`: Enable interactive mode with guided fixes
- `-q, --quiet`: Suppress output except for errors
- `-v, --version`: Display version information
- `-h, --help`: Show help information

## Examples

### Non-Interactive Mode

```bash
$ gitstatus
Current branch: main
⚠ You have uncommitted changes
```

### Interactive Mode

```bash
$ gitstatus -i
Current branch: main
⚠ You have uncommitted changes

=== Recommended Actions ===
1. Commit your changes:
git add .
git commit -m "your commit message"
Would you like to commit all changes? (y/n): y
Enter commit message: update readme
✓ Changes committed successfully

⚠ You have unpushed commits
2. Push your commits:
git push origin main
Would you like to push your commits? (y/n): y
✓ Commits pushed successfully

✓ All issues were resolved successfully
```

In interactive mode, the tool will:
1. Check for uncommitted changes and offer to commit them
2. Check for unpushed commits and offer to push them
3. Verify branch protection rules
4. Check if the branch needs to be updated

Each step is handled sequentially with clear prompts and confirmations.

