# AI Output Validator

「AI成果物のUnit Test」— comet-taskAI ロードマップ Product I。

acceptance_criteria（受け入れ条件）を定義したテストケースに対して、AIが生成した出力を
決定論的なルールで検証し、pass/fail・スコアを返す。CI/CDパイプラインに組み込める。

本製品自身はどのAIプロバイダーも呼び出さない。外部システム（CI・ai-scheduler・独自スクリプト等）
がAIを実行し、その出力をこの製品にPOSTして検証結果を受け取る、という構成。

## 現在のステータス: Phase 1-2（ルールエンジン・CRUD API・CLI）完了

- [x] Phase 0: プロジェクト立ち上げ
- [x] Phase 1: データモデル・ルールエンジン・CRUD API
- [x] Phase 2: CLI（`ovcli`、CI/CD組み込み用）
- [ ] Phase 3: Wails + Vue3 UI
- [ ] Phase 4: 仕上げ・署名・配布・LP

## 使い方（開発用ヘッドレスサーバー）

```bash
go mod tidy
go run ./cmd/ovserve   # :8426 でAPIサーバー起動
go run ./cmd/smoketest
```

### Suite・Caseの作成、実行

```bash
curl -X POST localhost:8426/api/v1/suites \
  -H "Authorization: Bearer $API_KEY" -H "Content-Type: application/json" \
  -d '{"name":"greeting-bot"}'

curl -X POST localhost:8426/api/v1/suites/{suiteId}/cases \
  -H "Authorization: Bearer $API_KEY" -H "Content-Type: application/json" \
  -d '{"name":"says hello","rules":[{"type":"contains","value":"Hello"}]}'

curl -X POST localhost:8426/api/v1/suites/{suiteId}/run \
  -H "Authorization: Bearer $API_KEY" -H "Content-Type: application/json" \
  -d '{"source":"ci","cases":[{"case_id":"{caseId}","output":"Hello there!"}]}'
```

### ルール種別

| type | valueの意味 |
|---|---|
| `contains` | 部分文字列（含む） |
| `not_contains` | 部分文字列（含まない） |
| `regex` | 正規表現パターン |
| `min_length` / `max_length` | 文字数 |
| `json_valid` | 出力が妥当なJSONか（valueは無視） |
| `json_key_exists` | JSONオブジェクトの特定トップレベルキーが存在するか |

1テストケース内の複数ルールはAND条件（全て合格して初めてケース合格）。

### CI/CDへの組み込み（`ovcli`）

```bash
export OV_API_KEY=...
echo '{"source":"ci:github-actions","cases":[{"case_id":"...","output":"..."}]}' \
  | go run ./cmd/ovcli -suite {suiteId}

# 不合格ケースが1件でもあれば exit 1 になるので、そのままCIのステップとして使える
```

## API

| メソッド | パス | 用途 |
|---|---|---|
| POST/GET/DELETE | `/api/v1/keys` | APIキー管理 |
| POST/GET/DELETE | `/api/v1/suites` | Suite CRUD |
| POST/GET/DELETE | `/api/v1/suites/{id}/cases` | Case CRUD |
| POST | `/api/v1/suites/{id}/run` | 検証実行（TestRun + CaseResultを作成） |
| GET | `/api/v1/suites/{id}/runs` | 実行履歴 |
| GET | `/api/v1/runs/{id}` | 実行詳細（ケース・ルールごとの内訳） |

## ディレクトリ構成

```
internal/db/       GORMモデル（TestSuite/TestCase/TestRun/CaseResult/AgentKey）
internal/api/       REST API（suites/cases/runs/keys）+ ルールエンジン + 認証ミドルウェア
cmd/ovserve/        開発用ヘッドレスAPIサーバー
cmd/ovcli/          CI/CD組み込み用CLI
cmd/smoketest/      通しスモークテスト
docs/                設計ドキュメント
```
