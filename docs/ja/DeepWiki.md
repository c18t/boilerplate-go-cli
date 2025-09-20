# boilerplate-go-cli DeepWiki (日本語版)

## 目次

1. [プロジェクト概要](#プロジェクト概要)
2. [アーキテクチャ](#アーキテクチャ)
3. [開発環境セットアップ](#開発環境セットアップ)
4. [依存性注入](#依存性注入)
5. [コマンド生成](#コマンド生成)
6. [ビルドとリリース](#ビルドとリリース)
7. [開発ワークフロー](#開発ワークフロー)
8. [ファイル構造](#ファイル構造)

## プロジェクト概要

`boilerplate-go-cli`は、Go言語でCLIツールを開発するためのテンプレートプロジェクトです。クリーンアーキテクチャパターンと依存性注入を採用し、保守性と拡張性に優れたCLIアプリケーションの開発を支援します。

### 主な特徴

- **クリーンアーキテクチャ**: ビジネスロジックとインフラストラクチャの分離
- **依存性注入**: `samber/do/v2`を使用した効率的な依存関係管理
- **コード生成**: `scaffdog`による定型コードの自動生成
- **開発環境**: Dev Containersによる一貫した開発環境
- **ツール管理**: `mise`による開発ツールのバージョン管理
- **自動リリース**: `GoReleaser`による多プラットフォーム対応

### 技術スタック

- **言語**: Go 1.22.5
- **CLIフレームワーク**: Cobra
- **依存性注入**: samber/do/v2
- **コード生成**: scaffdog
- **ツール管理**: mise
- **リリース**: GoReleaser
- **開発環境**: Dev Containers

## アーキテクチャ

このプロジェクトは、クリーンアーキテクチャパターンを採用しており、以下の層で構成されています：

### アーキテクチャ図

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Invoker       │───▶│   Controller    │───▶│  Use Case Bus   │
│  (cmd/*.go)     │    │ (controller/)   │    │   (port/)       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                        │
                                                        ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Presenter     │◀───│   Interactor    │◀───│    Use Case     │
│ (presenter/)    │    │ (interactor/)   │    │   Interface     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### 各層の責務

#### 1. Invoker層 (`cmd/`)
- Cobraコマンドの定義
- 依存性注入コンテナからControllerを取得
- コマンド実行の起点

#### 2. Controller層 (`internal/adapter/controller/`)
- HTTPリクエストやCLI引数の処理
- Use Case Busへのデータ転送
- 入力データの検証

#### 3. Use Case Bus層 (`internal/usecase/port/`)
- 複数のUse Caseを統合管理
- 入力データの型に基づく適切なUse Caseへのルーティング
- Use Caseインターフェースの定義

#### 4. Interactor層 (`internal/usecase/interactor/`)
- ビジネスロジックの実装
- ドメインルールの適用
- Presenterへの結果転送

#### 5. Presenter層 (`internal/adapter/presenter/`)
- 出力フォーマットの制御
- エラーハンドリング
- ユーザーインターフェースへの結果表示

### コア型定義

```go
// internal/core/type.go
type RunEFunc func(cmd *cobra.Command, args []string) (err error)

type Controller interface {
    Exec(cmd *cobra.Command, args []string) (err error)
}

type UseCase interface{}
```

## 開発環境セットアップ

### 前提条件

以下のいずれかの環境が必要です：

- [GitHub Codespaces](https://github.co.jp/features/codespaces)
- [Dev Containers](https://code.visualstudio.com/docs/devcontainers/containers)対応IDE（Visual Studio Codeなど）

### セットアップ手順

#### GitHub Codespacesを使用する場合

1. プロジェクトページで「Use this template」をクリック
2. 「Open in a codespace」を選択
3. Codespaces環境でセットアップを実行

#### Dev Containersを使用する場合

1. テンプレートから新しいリポジトリを作成：
   ```shell
   gh repo create my-project --template c18t/boilerplate-go-cli
   ```

2. プロジェクトをクローンして移動：
   ```shell
   ghq get <name>/my-project
   cd $(ghq root)/github.com/<name>/my-project
   ```

3. 環境変数を設定：
   ```shell
   cp .env.sample .env
   (echo UID=$(id -u) & echo GID=$(id -g)) >> .env
   ```

4. Dev Containersで開く：
   - `code .`
   - `Ctrl` + `Shift` + `P`
   - `>Dev Containers: Reopen in Container`

### コンテナワークスペースのセットアップ

1. セットアップタスクを実行：
   ```shell
   mise trust
   mise run setup
   ```

2. 新しいコマンドを作成：
   ```shell
   cobra-cli init
   cobra-cli add <new command>
   scaffdog generate command --answer "name:<new command>" --answer "usecase:command"
   ```

3. コマンドとコントローラーを接続：
   ```diff
   func init() {
   +   testCmd.RunE = createTestCommand()
       rootCmd.AddCommand(testCmd)
   ```

4. アプリケーションをビルドして実行：
   ```shell
   mise run build
   ./bin/app
   ```

## 依存性注入

このプロジェクトでは、`samber/do/v2`ライブラリを使用して依存性注入を実装しています。

### 基本構造

```go
// internal/inject/000_inject.go
package inject

import "github.com/samber/do/v2"

var Injector = AddProvider()

func AddProvider() *do.RootScope {
    var i = do.New()
    
    // 依存性の登録
    // do.Provide(i, gateway.SomeGateway)
    // do.Provide(i, repository.SomeRepository)
    
    return i
}
```

### コマンド固有の依存性注入

各コマンドには専用の依存性注入コンテナが作成されます：

```go
// internal/inject/{command_name}.go
var Injector{CommandName} = Add{CommandName}Provider()

func Add{CommandName}Provider() *do.RootScope {
    var i = Injector.Clone()
    
    // adapter/controller
    do.Provide(i, controller.New{CommandName}Controller)
    
    // usecase/port
    do.Provide(i, port.New{CommandName}UseCaseBus)
    
    // usecase/interactor
    do.Provide(i, interactor.New{CommandName}Interactor)
    
    // adapter/presenter
    do.Provide(i, presenter.New{CommandName}Presenter)
    
    return i
}
```

### 依存性の解決

Invokerでコントローラーを取得：

```go
// cmd/{command_name}_invoker.go
func create{CommandName}Command() core.RunEFunc {
    cmd, err := do.Invoke[controller.{CommandName}Controller](inject.Injector{CommandName})
    cobra.CheckErr(err)
    return cmd.Exec
}
```

## コマンド生成

`scaffdog`を使用して、新しいコマンドの定型コードを自動生成できます。

### 生成されるファイル

1. **Invoker** (`cmd/{command_name}_invoker.go`)
2. **Controller** (`internal/adapter/controller/{command_name}.go`)
3. **Use Case Port** (`internal/usecase/port/{command_name}.go`)
4. **Interactor** (`internal/usecase/interactor/{command_name}.go`)
5. **Presenter** (`internal/adapter/presenter/{command_name}.go`)
6. **依存性注入設定** (`internal/inject/{command_name}.go`)

### 生成コマンド

```shell
scaffdog generate command --answer "name:<command_name>" --answer "usecase:command"
```

### 生成されるコードの例

#### Controller

```go
type {CommandName}Controller interface {
    core.Controller
    Params() *{CommandName}Params
}

type {commandName}Controller struct {
    bus    port.{CommandName}UseCaseBus `do:""`
    params *{CommandName}Params
}

func (c *{commandName}Controller) Exec(cmd *cobra.Command, args []string) (err error) {
    c.bus.Handle(&port.{CommandName}UseCaseInputData{})
    return
}
```

#### Use Case Bus

```go
type {CommandName}UseCaseBus interface {
    Handle(input {CommandName}UseCaseInputData)
}

func (bus *{commandName}UseCaseBus) Handle(input {CommandName}UseCaseInputData) {
    switch data := input.(type) {
    case *{CommandName}UseCaseInputData:
        bus.command.Handle(data)
    default:
        panic(fmt.Errorf("handler for '%T' is not implemented", data))
    }
}
```

## ビルドとリリース

### ビルドシステム

#### Makefile

```makefile
# バージョン情報
VER_TAG:=$(shell git describe --tag)
VER_REV:=$(shell git rev-parse --short HEAD)
VERSION:=${VER_TAG}+${VER_REV}

# ビルドフラグ
GO_LDFLAGS:=-s -w -extldflags "-static" -X "main.version=${VERSION}"
GO_BUILD:=-ldflags '${GO_LDFLAGS}' -trimpath

# ビルドターゲット
build: app

app: main.go
	go build -o $(BINDIR)/$(notdir $(shell pwd)) $(GO_BUILD) $^
```

#### mise タスク

```toml
[tasks.build]
description = "Build the CLI application"
alias = "b"
run = "make"
sources = ["go.mod", "main.go", "cmd/**/*.go", "internal/**/*.go"]
outputs = ["bin/app"]

[tasks.release]
description = "Build release binaries"
alias = "r"
run = "goreleaser release --snapshot --clean"
```

### リリースプロセス

#### GoReleaser設定

```yaml
# .goreleaser.yaml
builds:
  - env:
      - CGO_ENABLED=0
    goos:
      - linux
      - windows
      - darwin

archives:
  - format: tar.gz
    name_template: >-
      {{ .ProjectName }}_
      {{- title .Os }}_
      {{- if eq .Arch "amd64" }}x86_64
      {{- else if eq .Arch "386" }}i386
      {{- else }}{{ .Arch }}{{ end }}
```

#### 自動リリースの有効化

1. サンプルワークフローを有効化：
   ```shell
   mv .github/workflows/release.yaml.example .github/workflows/release.yaml
   ```

2. `main`ブランチにプッシュ

3. [tagpr](https://github.com/Songmu/tagpr)ボットからのプルリクエストを承認

4. リリースページでリリースを確認

## 開発ワークフロー

### 利用可能なタスク

`mise run <task name>`で実行できるタスク：

| タスク名 | 説明 | エイリアス |
|---------|------|-----------|
| `setup` | 全セットアップタスクの実行 | - |
| `setup:mise` | miseによる開発依存関係のインストール | - |
| `setup:go-mod` | go.modによるGoモジュールのインストール | - |
| `setup:pre-commit` | pre-commitフックのセットアップ | - |
| `build` | CLIアプリケーションのビルド | `b` |
| `release` | リリースバイナリのビルド | `r` |

### 開発ツール

#### mise設定

```toml
[tools]
node = "22.5.1"
goreleaser = "2.1.0"
pre-commit = "3.7.1"
shellcheck = "0.10.0"
"go:github.com/spf13/cobra-cli" = "1.3.0"
"npm:prettier" = "3.3.3"
"npm:scaffdog" = "4.0.0"
```

#### 推奨VS Code拡張機能

プロジェクトには以下の拡張機能が推奨されています：

- Go言語サポート
- YAML言語サポート
- GitHub Actions
- GitLens
- Prettier

### デバッグ環境

Dev Containersには、Delveデバッガーが組み込まれています：

```yaml
# compose.yaml
services:
  app:
    ports:
      - ${DEBUG_PORT:-2345}:${DEBUG_PORT:-2345}
    environment:
      - DELVE_VERSION=${DELVE_VERSION:-v1.22.1}
      - DEBUG_PORT=${DEBUG_PORT:-2345}
```

## ファイル構造

```
boilerplate-go-cli/
├── .devcontainer/          # Dev Container設定
│   ├── devcontainer.json
│   ├── Dockerfile
│   └── .cobra.yaml
├── .github/                # GitHub設定
│   ├── workflows/
│   │   └── release.yaml.example
│   ├── dependabot.yaml
│   └── pull_request_template.md
├── .scaffdog/              # コード生成テンプレート
│   ├── config.js
│   └── command.md
├── .vscode/                # VS Code設定
│   ├── extensions.json
│   ├── launch.json
│   └── settings.json
├── internal/               # 内部パッケージ
│   ├── core/              # コアインターフェース
│   │   └── type.go
│   └── inject/            # 依存性注入
│       └── 000_inject.go
├── docs/                   # ドキュメント
│   └── ja/                # 日本語ドキュメント
│       └── DeepWiki.md    # このファイル
├── .env.sample            # 環境変数サンプル
├── .goreleaser.yaml       # リリース設定
├── .mise.toml             # ツール管理設定
├── compose.yaml           # Docker Compose設定
├── go.mod                 # Goモジュール定義
├── go.sum                 # 依存関係チェックサム
├── Makefile               # ビルド設定
└── README.md              # プロジェクト概要
```

### 生成されるファイル（コマンド作成後）

```
├── cmd/                    # コマンド定義
│   ├── root.go
│   ├── {command}.go
│   └── {command}_invoker.go
├── internal/
│   ├── adapter/
│   │   ├── controller/
│   │   │   └── {command}.go
│   │   └── presenter/
│   │       └── {command}.go
│   ├── usecase/
│   │   ├── interactor/
│   │   │   └── {command}.go
│   │   └── port/
│   │       └── {command}.go
│   └── inject/
│       └── {command}.go
├── bin/                    # ビルド成果物
│   └── app
└── main.go                 # エントリーポイント
```

---

このDeepWikiは、`boilerplate-go-cli`プロジェクトの包括的なガイドです。新しい機能の追加や既存コードの理解に活用してください。

質問や改善提案がある場合は、GitHubのIssueまたはPull Requestでお知らせください。
