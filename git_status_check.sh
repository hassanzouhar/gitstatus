#!/bin/bash

# Script information
VERSION="1.0.0"
LAST_UPDATED="2024-01-10"

# Color codes
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default settings
QUIET_MODE=false
INTERACTIVE_MODE=false
PROTECTED_BRANCHES=("main" "master" "production" "staging")
EXIT_CODE=0

# Help message
show_help() {
    cat << EOF
Git Status Check v${VERSION}
Last updated: ${LAST_UPDATED}

Usage: $(basename $0) [-h] [-q] [-v] [-i]

Options:
-h    Show this help message
-q    Quiet mode (only return exit code, useful for CI/CD)
-v    Show version information
-i    Interactive mode (execute recommended actions with confirmation)

Exit codes:
0    Repository is clean and up-to-date
1    Not a git repository
2    Uncommitted changes exist
3    Unpushed/unpulled changes exist
4    Branches have diverged
5    Protected branch warning
6    Stashed changes exist
EOF
    exit 0
}

# Version information
show_version() {
    echo "Git Status Check v${VERSION}"
    echo "Last updated: ${LAST_UPDATED}"
    exit 0
}

# Parse command line options
# Function to execute git commands safely
execute_git_command() {
    local cmd="$1"
    local success_msg="$2"
    local error_msg="$3"
    
    if [ "$QUIET_MODE" = false ]; then
        echo -e "${BLUE}Executing: ${NC}$cmd"
    fi
    
    if eval "$cmd"; then
        [ "$QUIET_MODE" = false ] && echo -e "${GREEN}✓ $success_msg${NC}"
        return 0
    else
        [ "$QUIET_MODE" = false ] && echo -e "${RED}✗ $error_msg${NC}"
        return 1
    fi
}

# Function to prompt user for confirmation
confirm_action() {
    local prompt="$1"
    if [ "$INTERACTIVE_MODE" = false ] || [ "$QUIET_MODE" = true ]; then
        return 1
    fi
    
    while true; do
        echo -en "${YELLOW}$prompt ${NC}(y/n): "
        read -r response
        case "$response" in
            [yY][eE][sS]|[yY]) 
                return 0
                ;;
            [nN][oO]|[nN])
                return 1
                ;;
            *)
                echo -e "${RED}Please answer yes/y or no/n${NC}"
                ;;
        esac
    done
}

while getopts ":hqvi" opt; do
    case ${opt} in
        h )
            show_help
            ;;
        q )
            QUIET_MODE=true
            ;;
        v )
            show_version
            ;;
        i )
            INTERACTIVE_MODE=true
            ;;
        \? )
            echo -e "${RED}Error: Invalid option: -$OPTARG${NC}" >&2
            show_help
            exit 1
            ;;
        : )
            echo -e "${RED}Error: Option -$OPTARG requires an argument${NC}" >&2
            show_help
            exit 1
            ;;
    esac
done
shift $((OPTIND -1))

# Function to check if Git LFS is installed and being used
check_git_lfs() {
    if command -v git-lfs >/dev/null 2>&1; then
        if [ -f .gitattributes ] && grep -q "filter=lfs" .gitattributes; then
            if ! git lfs status >/dev/null 2>&1; then
                [ "$QUIET_MODE" = false ] && echo -e "${YELLOW}⚠ Git LFS is configured but not initialized${NC}"
                return 1
            fi
        fi
    fi
    return 0
}

# Function to check for stashed changes
check_stashed_changes() {
    if [ -n "$(git stash list)" ]; then
        [ "$QUIET_MODE" = false ] && echo -e "${YELLOW}⚠ You have stashed changes:${NC}"
        [ "$QUIET_MODE" = false ] && git stash list
        return 1
    fi
    return 0
}

# Function to check if current branch is protected
check_protected_branch() {
    local current_branch="$1"
    for branch in "${PROTECTED_BRANCHES[@]}"; do
        if [ "$current_branch" = "$branch" ]; then
            [ "$QUIET_MODE" = false ] && echo -e "${YELLOW}⚠ Warning: You are on protected branch '$branch'${NC}"
            return 1
        fi
    done
    return 0
}

# Initialize state variables
has_uncommitted_changes=false
has_unpushed_commits=false
needs_pull=false
has_diverged=false
# Check if current directory is a git repository
if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    [ "$QUIET_MODE" = false ] && echo -e "${RED}Error: Not a git repository${NC}"
    exit 1
fi

# Check Git LFS status
check_git_lfs || EXIT_CODE=2

# Check for stashed changes
check_stashed_changes || EXIT_CODE=6

# Get current branch name
current_branch=$(git symbolic-ref --short HEAD 2>/dev/null)
[ "$QUIET_MODE" = false ] && echo -e "${GREEN}Current branch: ${NC}$current_branch"

# Check if current branch is protected
check_protected_branch "$current_branch" || EXIT_CODE=5

# Check for uncommitted changes
if [[ -n $(git status --porcelain) ]]; then
    echo -e "${YELLOW}⚠ You have uncommitted changes:${NC}"
    git status -s
    has_uncommitted_changes=true
else
    echo -e "${GREEN}✓ Working directory is clean${NC}"
fi

# Check for unpushed commits
if [ -n "$(git log @{u}.. 2>/dev/null)" ]; then
    echo -e "${YELLOW}⚠ You have unpushed commits${NC}"
    git log @{u}.. --oneline
    has_unpushed_commits=true
else
    echo -e "${GREEN}✓ All commits are pushed${NC}"
fi

# Check for unpulled commits
git fetch -q
UPSTREAM=${1:-'@{u}'}
if ! git rev-parse --verify @{u} >/dev/null 2>&1; then
    echo -e "${YELLOW}⚠ No upstream branch set${NC}"
    has_unpushed_commits=true
else
    LOCAL=$(git rev-parse @ 2>/dev/null)
    REMOTE=$(git rev-parse @{u} 2>/dev/null)
    BASE=$(git merge-base @ @{u} 2>/dev/null)

    if [ "$LOCAL" = "$REMOTE" ]; then
        echo -e "${GREEN}✓ Repository is up to date${NC}"
    elif [ "$LOCAL" = "$BASE" ]; then
        BEHIND_COUNT=$(git rev-list --count HEAD..@{u} 2>/dev/null)
        echo -e "${YELLOW}⚠ Need to pull - you are behind by ${BEHIND_COUNT} commit(s)${NC}"
        needs_pull=true
    elif [ "$REMOTE" = "$BASE" ]; then
        AHEAD_COUNT=$(git rev-list --count @{u}..HEAD 2>/dev/null)
        echo -e "${YELLOW}⚠ Need to push - you are ahead by ${AHEAD_COUNT} commit(s)${NC}"
        has_unpushed_commits=true
    else
        echo -e "${RED}✗ Branches have diverged${NC}"
        has_diverged=true
    fi
fi

# Print summary and recommendations
echo -e "\n${BLUE}=== Recommended Actions ===${NC}"
if [ "$has_uncommitted_changes" = true ]; then
    echo -e "${YELLOW}1. Commit your changes:${NC}"
    echo "   git add ."
    echo "   git commit -m \"your commit message\""
    
    if confirm_action "Would you like to commit all changes?"; then
        read -p "Enter commit message: " commit_msg
        execute_git_command "git add . && git commit -m \"$commit_msg\"" \
            "Changes committed successfully" \
            "Failed to commit changes"
    fi
fi

if [ "$has_unpushed_commits" = true ]; then
    echo -e "${YELLOW}2. Push your commits:${NC}"
    echo "   git push origin $current_branch"
    
    if confirm_action "Would you like to push your commits?"; then
        execute_git_command "git push origin $current_branch" \
            "Commits pushed successfully" \
            "Failed to push commits"
    fi
fi

if [ "$needs_pull" = true ]; then
    echo -e "${YELLOW}3. Pull latest changes:${NC}"
    echo "   git pull origin $current_branch"
    
    if confirm_action "Would you like to pull latest changes?"; then
        execute_git_command "git pull origin $current_branch" \
            "Changes pulled successfully" \
            "Failed to pull changes"
    fi
fi

if [ "$has_diverged" = true ]; then
    echo -e "${YELLOW}4. Resolve diverged branches:${NC}"
    echo "   git pull origin $current_branch  # Pull and merge remote changes"
    echo "   # Resolve any conflicts if they occur"
    echo "   git push origin $current_branch  # Push your changes"
    
    if confirm_action "Would you like to sync diverged branches?"; then
        if execute_git_command "git pull origin $current_branch" \
            "Remote changes pulled successfully" \
            "Failed to pull remote changes"; then
            echo -e "${YELLOW}Check for conflicts and resolve if necessary${NC}"
            if confirm_action "Ready to push changes?"; then
                execute_git_command "git push origin $current_branch" \
                    "Changes pushed successfully" \
                    "Failed to push changes"
            fi
        fi
    fi
fi

if [[ "$has_uncommitted_changes" = false && "$has_unpushed_commits" = false && "$needs_pull" = false && "$has_diverged" = false ]]; then
    [ "$QUIET_MODE" = false ] && echo -e "${GREEN}✓ No actions needed - repository is in sync${NC}"
fi

# Set final exit code based on repository state
if [ "$has_uncommitted_changes" = true ]; then
    exit 2
elif [ "$has_unpushed_commits" = true ] || [ "$needs_pull" = true ]; then
    exit 3
elif [ "$has_diverged" = true ]; then
    exit 4
elif [ $EXIT_CODE -ne 0 ]; then
    exit $EXIT_CODE
else
    exit 0
fi
