// cmd/ovcli はCI/CDパイプラインに組み込むためのコマンドラインクライアント。
// 標準入力（またはファイル）から {"source":"...","cases":[{"case_id":"...","output":"..."}]}
// 形式のJSONを読み、サーバーの /api/v1/suites/:id/run にPOSTする。
// 全ケース合格ならexit 0、1件でも不合格ならexit 1、実行時エラーならexit 2。
//
//	go run ./cmd/ovcli -addr http://127.0.0.1:8426 -key $OV_API_KEY -suite $SUITE_ID -file result.json
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

type runCaseInput struct {
	CaseID string `json:"case_id"`
	Output string `json:"output"`
}

type runSuiteInput struct {
	Source string         `json:"source"`
	Cases  []runCaseInput `json:"cases"`
}

type ruleResult struct {
	Type    string `json:"type"`
	Value   string `json:"value"`
	Passed  bool   `json:"passed"`
	Message string `json:"message"`
}

type caseResult struct {
	CaseName    string       `json:"case_name"`
	Passed      bool         `json:"passed"`
	RuleResults []ruleResult `json:"rule_results"`
}

type testRun struct {
	Passed bool    `json:"passed"`
	Score  float64 `json:"score"`
}

type runResult struct {
	Run     testRun      `json:"run"`
	Results []caseResult `json:"results"`
}

func main() {
	addr := flag.String("addr", envOr("OV_ADDR", "http://127.0.0.1:8426"), "AI Output Validator API base URL")
	key := flag.String("key", os.Getenv("OV_API_KEY"), "APIキー（またはOV_API_KEY環境変数）")
	suite := flag.String("suite", "", "実行対象のSuite ID（必須）")
	file := flag.String("file", "", "入力JSONファイル（省略時は標準入力）")
	flag.Parse()

	if *suite == "" {
		fmt.Fprintln(os.Stderr, "error: -suite is required")
		os.Exit(2)
	}
	if *key == "" {
		fmt.Fprintln(os.Stderr, "error: -key or OV_API_KEY is required")
		os.Exit(2)
	}

	var body []byte
	var err error
	if *file != "" {
		body, err = os.ReadFile(*file)
	} else {
		body, err = io.ReadAll(os.Stdin)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading input: %s\n", err)
		os.Exit(2)
	}

	var in runSuiteInput
	if err := json.Unmarshal(body, &in); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing input JSON: %s\n", err)
		os.Exit(2)
	}

	req, err := http.NewRequest(http.MethodPost, *addr+"/api/v1/suites/"+*suite+"/run", bytes.NewReader(body))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error building request: %s\n", err)
		os.Exit(2)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+*key)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error calling server: %s\n", err)
		os.Exit(2)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading response: %s\n", err)
		os.Exit(2)
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "server returned %d: %s\n", resp.StatusCode, string(respBody))
		os.Exit(2)
	}

	var result runResult
	if err := json.Unmarshal(respBody, &result); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing server response: %s\n", err)
		os.Exit(2)
	}

	printSummary(result)

	if !result.Run.Passed {
		os.Exit(1)
	}
}

func printSummary(result runResult) {
	status := "PASS"
	if !result.Run.Passed {
		status = "FAIL"
	}
	fmt.Printf("%s  score=%.2f  (%d cases)\n\n", status, result.Run.Score, len(result.Results))

	for _, c := range result.Results {
		mark := "✓"
		if !c.Passed {
			mark = "✗"
		}
		fmt.Printf("%s %s\n", mark, c.CaseName)
		for _, rr := range c.RuleResults {
			if rr.Passed {
				continue
			}
			fmt.Printf("    - [%s] %s\n", rr.Type, rr.Message)
		}
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
