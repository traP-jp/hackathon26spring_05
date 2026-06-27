# Qpid backend file structure

Qpid は Go + Echo の backend です。`main.go` から `handler.Serve()` を呼び出し、HTTP ルーティング、repository、domain model を分けています。

## Entry point

- `main.go`
  - backend のエントリーポイントです。
  - 現在は `handler.Serve()` を呼ぶだけです。

## HTTP layer

- `handler/`
  - Echo の handler とルーティングを置きます。
  - `handler.go`
    - `Serve()` で Echo を作成し、repository、session store、middleware を組み立てます。
    - `mapRoutes()` で `/api` 以下の route を定義します。
  - `error.go`
    - handler から返す error response 型を置きます。
- `handler/middleware/`
  - HTTP middleware と handler から使う request-scoped helper を置きます。

## Domain layer

- `domain/`
  - repository や handler から共有する domain model と小さい interface を置きます。

## Repository layer

- `repository/`
  - backend が必要とする永続化操作の interface を置きます。
- `repository/mock/`
  - repository interface の mock 実装を置きます。

## Infrastructure layer

- `infrastructure/`
  - DB 接続や DB を使う実 repository など、外部インフラ依存の実装を置きます。

## Database and runtime files

- `db/init/`
  - Docker Compose で DB 初期化 SQL を置く場所です。
- `docker-compose.yml`
  - backend 開発用の container 構成です。
- `dev.Dockerfile`, `prod.Dockerfile`
  - 開発用、本番用の backend image 定義です。
- `Taskfile.yml`
  - 開発用 task 定義です。
- `.env.example`
  - backend 起動に必要な環境変数の例です。

## When adding backend code

- HTTP endpoint を追加するときは `handler/handler.go` の `mapRoutes()` に route を足し、処理本体は用途別の `handler/*.go` に置きます。
- 永続化操作が必要なときは、まず `repository/` に interface method を追加し、`repository/mock/` と `infrastructure/` 側の実装を合わせます。
- handler と repository の間で共有する data shape は `domain/` に置きます。
