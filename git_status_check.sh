#!/bin/bash

# Color codes
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Initialize state variables
has_uncommitted_changes=false
has_unpushed_commits=false
needs_pull=false
has_diverged=false
# Check if current directory is a git repository
if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    echo -e "${RED}Error: Not a git repository${NC}"
    exit 1
fi

# Get current branch name
current_branch=$(git symbolic-ref --short HEAD 2>/dev/null)
echo -e "${GREEN}Current branch: ${NC}$current_branch"

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
LOCAL=$(git rev-parse @)
REMOTE=$(git rev-parse "$UPSTREAM")
BASE=$(git merge-base @ "$UPSTREAM")

if [ $LOCAL = $REMOTE ]; then
    echo -e "${GREEN}✓ Repository is up to date${NC}"
elif [ $LOCAL = $BASE ]; then
    echo -e "${YELLOW}⚠ Need to pull - you are behind by $(git rev-list HEAD..origin/$(git branch --show-current) --count) commit(s)${NC}"
    needs_pull=true
elif [ $REMOTE = $BASE ]; then
    echo -e "${YELLOW}⚠ Need to push - you are ahead by $(git rev-list origin/$(git branch --show-current)..HEAD --count) commit(s)${NC}"
    has_unpushed_commits=true
else
    echo -e "${RED}✗ Branches have diverged${NC}"
    has_diverged=true
fi

# Print summary and recommendations
echo -e "\n${BLUE}=== Recommended Actions ===${NC}"
if [ "$has_uncommitted_changes" = true ]; then
    echo -e "${YELLOW}1. Commit your changes:${NC}"
    echo "   git add ."
    echo "   git commit -m \"your commit message\""
fi

if [ "$has_unpushed_commits" = true ]; then
    echo -e "${YELLOW}2. Push your commits:${NC}"
    echo "   git push origin $current_branch"
fi

if [ "$needs_pull" = true ]; then
    echo -e "${YELLOW}3. Pull latest changes:${NC}"
    echo "   git pull origin $current_branch"
fi

if [ "$has_diverged" = true ]; then
    echo -e "${YELLOW}4. Resolve diverged branches:${NC}"
    echo "   git pull origin $current_branch  # Pull and merge remote changes"
    echo "   # Resolve any conflicts if they occur"
    echo "   git push origin $current_branch  # Push your changes"
fi

if [[ "$has_uncommitted_changes" = false && "$has_unpushed_commits" = false && "$needs_pull" = false && "$has_diverged" = false ]]; then
    echo -e "${GREEN}✓ No actions needed - repository is in sync${NC}"
fi
