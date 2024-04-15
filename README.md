# boilerplate-go-cli

my template project for go-cli.

## Getting Started

### Prerequisites

- [Codespaces](https://github.co.jp/features/codespaces)

Or some IDE with [Dev Containers](https://code.visualstudio.com/docs/devcontainers/containers)
support (e.g., Visual Studio Code).

### Create Container

#### Using GitHub Codespaces

1. Click "Use this template" on the project page and select "Open in a codespace".
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
3. Open the project in Dev Containers:
   - `code .`
   - `Ctrl` + `Shift` + `P`
   - `>Dev Containers: Reopen in Container`

### Setup Container Workspace

1. Run setup tasks:
   ```shell
   mise trust
   mise run setup
   ```
2. Create new command:
   ```shell
   cobra-cli init
   cobra-cli add <new command>
   ```
3. Build and run the application:
   ```shell
   mise run build
   ./bin/app
   ```

## Avaiable Task Runner Commands

`mise run <task name>`

```console
$ mise tasks
Name              Description                         Source
build             Build the CLI application           ~/Projects/github.com/c18t/boilerplate-go-cli/.mise.toml
release           Build release binaries              ~/Projects/github.com/c18t/boilerplate-go-cli/.mise.toml
setup             Setup (Runs all `setup:*` tasks)    ~/Projects/github.com/c18t/boilerplate-go-cli/.mise.toml
setup:go-mod      Install go modules with go.mod      ~/Projects/github.com/c18t/boilerplate-go-cli/.mise.toml
setup:mise        Install dev dependencies with mise  ~/Projects/github.com/c18t/boilerplate-go-cli/.mise.toml
setup:pre-commit  Sets up pre-commit hooks            ~/Projects/github.com/c18t/boilerplate-go-cli/.mise.toml
```
