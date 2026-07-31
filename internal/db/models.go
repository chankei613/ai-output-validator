// Package db はAI Output ValidatorのGORMモデルとSQLite初期化を提供する。
package db

import "time"

// TestSuite はテストケースの集合（1つのプロンプト/機能に対応する）。
type TestSuite struct {
	ID          string `gorm:"primaryKey" json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Rule は1つの受け入れ条件。Typeに応じてValueの意味が変わる：
//   - contains/not_contains: 部分文字列
//   - regex: 正規表現パターン
//   - min_length/max_length: 文字数（数値文字列）
//   - json_valid: 未使用（Valueは空でよい）
//   - json_key_exists: JSONのトップレベルキー名
type Rule struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// TestCase は1つのテストケース（複数Ruleの組み合わせ = AND条件で全て満たせば合格）。
type TestCase struct {
	ID      string `gorm:"primaryKey" json:"id"`
	SuiteID string `gorm:"index" json:"suite_id"`
	Name    string `json:"name"`
	Rules   []Rule `gorm:"serializer:json" json:"rules"`

	CreatedAt time.Time `json:"created_at"`
}

// TestRun は「このSuiteを、この時点で、外部から届いた出力群に対して実行した」記録（追記専用）。
type TestRun struct {
	ID      string  `gorm:"primaryKey" json:"id"`
	SuiteID string  `gorm:"index" json:"suite_id"`
	Source  string  `json:"source"` // 呼び出し元の自由記述（例: "ci:github-actions"）
	Passed  bool    `json:"passed"`
	Score   float64 `json:"score"` // 合格ケース数 / 全ケース数

	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`
}

// RuleResult は1ルールぶんの評価結果。
type RuleResult struct {
	Type    string `json:"type"`
	Value   string `json:"value"`
	Passed  bool   `json:"passed"`
	Message string `json:"message"`
}

// CaseResult は1テストケースぶんの結果（ルールごとの内訳を含む）。
type CaseResult struct {
	ID          string       `gorm:"primaryKey" json:"id"`
	TestRunID   string       `gorm:"index" json:"test_run_id"`
	TestCaseID  string       `json:"test_case_id"`
	CaseName    string       `json:"case_name"`
	Output      string       `json:"output"`
	Passed      bool         `json:"passed"`
	RuleResults []RuleResult `gorm:"serializer:json" json:"rule_results"`
}

// AgentKey — CRUD/実行APIを叩くためのAPIキー。ハッシュのみ保存する。
type AgentKey struct {
	ID         string     `gorm:"primaryKey" json:"id"`
	Name       string     `json:"name"`
	APIKeyHash string     `json:"-"`
	CreatedAt  time.Time  `json:"created_at"`
	RevokedAt  *time.Time `json:"revoked_at"`
}
