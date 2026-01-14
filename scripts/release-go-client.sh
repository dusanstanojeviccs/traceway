#!/bin/bash
set -e

# Release script for go.tracewayapp.com Go client library

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
GO_CLIENT_DIR="$PROJECT_ROOT/clients/go-client"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

usage() {
    echo "Usage: $0 <version>"
    echo ""
    echo "Arguments:"
    echo "  version    Semantic version (e.g., v1.0.0, v0.1.0)"
    echo ""
    echo "Examples:"
    echo "  $0 v1.0.0"
    echo "  $0 v0.2.1"
    exit 1
}

# Check for version argument
if [ -z "$1" ]; then
    echo -e "${RED}Error: Version argument required${NC}"
    usage
fi

VERSION="$1"

# Validate semver format (vX.Y.Z)
if ! [[ "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo -e "${RED}Error: Invalid version format '$VERSION'${NC}"
    echo "Version must be in semver format: vX.Y.Z (e.g., v1.0.0)"
    exit 1
fi

echo -e "${GREEN}=== Releasing go.tracewayapp.com $VERSION ===${NC}"
echo ""

# Change to go-client directory
cd "$GO_CLIENT_DIR"
echo -e "${YELLOW}Working directory: $GO_CLIENT_DIR${NC}"
echo ""

# Verify the build
echo -e "${YELLOW}Building library...${NC}"
go build ./...
echo -e "${GREEN}Build successful${NC}"
echo ""

# Run vet
echo -e "${YELLOW}Running go vet...${NC}"
go vet ./...
echo -e "${GREEN}Vet passed${NC}"
echo ""

# Check for uncommitted changes
cd "$PROJECT_ROOT"
if [ -n "$(git status --porcelain)" ]; then
    echo -e "${RED}Error: You have uncommitted changes${NC}"
    echo "Please commit all changes before releasing."
    git status --short
    exit 1
fi

# Create the tag
TAG="clients/go-client/$VERSION"
echo -e "${YELLOW}Creating tag: $TAG${NC}"

git tag -a "$TAG" -m "Release go.tracewayapp.com $VERSION"
echo -e "${GREEN}Tag created${NC}"
echo ""

# Push the tag
echo -e "${YELLOW}Pushing tag to origin...${NC}"
git push origin "$TAG"
echo -e "${GREEN}Tag pushed${NC}"
echo ""

echo -e "${GREEN}=== Release Complete ===${NC}"
echo ""
echo "Next steps:"
echo "1. The Go module proxy will automatically pick up the new version"
echo "2. Users can install with: go get go.tracewayapp.com@$VERSION"
echo ""
echo "Note: Make sure go.tracewayapp.com is configured to redirect to your Git repository."
echo "See: https://go.dev/ref/mod#vcs-find for vanity import setup"
