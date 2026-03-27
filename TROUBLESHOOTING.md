# Troubleshooting

## Running Locally

You can run the action locally as a standalone Go binary for debugging or testing purposes.

### Build the binary

```bash
# macOS (Apple Silicon)
make build-darwin

# Linux
make build-linux

# Windows
make build-windows
```

### Set environment variables

The action reads inputs from environment variables using the `INPUT_` prefix:

```bash
export GITHUB_SHA="$(git rev-parse HEAD)"
export INPUT_BRANCHING_MODEL="git-flow"    # or "trunk-based"
export INPUT_PREFIX="v"
export INPUT_MAIN_BRANCH_NAME="master"
export INPUT_DEVELOP_BRANCH_NAME="develop"
export INPUT_REPO_DIR="."
export INPUT_DEBUG="true"
```

### Run the binary

```bash
# macOS
./build/darwin/arm64/semver

# Linux
./build/linux/amd64/semver

# Windows
./build/windows/amd64/semver.exe
```

### One-liner example

```bash
GITHUB_SHA="$(git rev-parse HEAD)" \
INPUT_BRANCHING_MODEL="trunk-based" \
INPUT_PREFIX="v" \
INPUT_MAIN_BRANCH_NAME="master" \
INPUT_REPO_DIR="." \
INPUT_DEBUG="true" \
./build/darwin/arm64/semver
```
