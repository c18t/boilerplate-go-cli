# boilerplate-go-cli

my template project for go-cli.

## Getting Started

### Prerequisites

- [Codespaces](https://github.co.jp/features/codespaces)

Or some IDE with [Dev Containers](https://code.visualstudio.com/docs/devcontainers/containers)
support (e.g., Visual Studio Code).

### Create Container

#### Using GitHub Codespaces

1. Click "Use this template" on this project page and select "Open in a codespace".
2. Follow the steps in the Codespaces environment to set up your development environment.
3. Once set up, proceed to the "Setup Container Workspace" section below.

#### Using Dev Containers

1. Create new repository using this template:
   ```shell
   gh repo create my-project --template c18t/boilerplate-go-cli
   ```
2. Clone and Navigate to the project directory:
   ```shell
   ghq get <name>/my-project
   cd $(ghq root)/github.com/<name>/my-project
   ```
3. Add GH_TOKEN to .env (if necessary):
   ```shell
   cp .env.sample .env
   gh auth token | xargs -I {} echo "GH_TOKEN="{} >> .env
   ```
4. Open the project in Dev Containers:
   1. `code .`
   1. `Ctrl` + `Shift` + `P`
   1. `>Dev Containers: Reopen in Container`

### Setup Container Workspace

1. Run setup tasks:
   ```shell
   post-create.sh
   ```
2. Create new command:
   ```shell
   cobra-cli init
   cobra-cli add <new command>
   scaffdog generate command --answer "name:<new command>" --answer "usecase:command"
   ```
3. Wire a command and a controller: (open `<new command>` code [e.g. ./cmd/test.go])
   ```diff
   func init() {
   +   testCmd.RunE = createTestCommand()
       rootCmd.AddCommand(testCmd)
   ```
4. Build and run the application:
   ```shell
   mise run build
   ./bin/app
   ```
5. [extra] Install extensions recommended for the workspace:
   1. `Ctrl` + `Shift` + `P`
   1. `>Extensions: Show Recommended Extensions`
   1. Click `install` button.

## Available Task Runner Commands

`mise run <task name>`

```console
$ mise tasks
Name                                 Description
build                                Build the CLI application
dev-up:ccmanager-skip-permissions    Set up ccmanager to skip permissions
dev-up:ccmanager-worktree-settings   Set up ccmanager worktree auto-directory settings
dev-up:claude-code-stop-autoupdates  Set up Claude Code to disable auto-updates
devcontainer-up                      Start devcontainer and run ccmanager with Claude...
release                              Build release binaries
setup                                Set up (Runs all `setup:*` tasks)
setup:claude-mcp                     Set up Claude Code MCP servers
setup:go-mod                         Install go modules with go.mod
setup:ignore-workspace-file-changes  Ignore local changes to workspace file
setup:mise                           Install dev dependencies with mise
setup:pnpm                           Set up pnpm packages
setup:pre-commit                     Set up pre-commit hooks
```

## Enabling Automated Releases

1. Enable the Sample Workflow:
   ```shell
   mv .github/workflows/release.yaml.example .github/workflows/release.yaml
   ```
2. Push to the `main` branch
3. Approve the pull request from the [tagpr](https://github.com/Songmu/tagpr) bot
4. Check the releases page of your repository
