// cmd/smoketest はAI Output ValidatorのAPIを一時DBで自前起動し、
// ブートストラップ鍵発行 → Suite/Case作成 → 実行（合格ケース）→ 実行（不合格ケース）→
// 実行履歴・詳細取得 → 未知case_idのエラー処理、の一連が通しで動くことを確認する。
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/chankei613/ai-output-validator/internal/api"
	"github.com/chankei613/ai-output-validator/internal/db"
)

func main() {
	dbPath := "smoketest.db"
	_ = os.Remove(dbPath)
	defer func() { _ = os.Remove(dbPath) }()

	conn, err := db.Init(dbPath)
	if err != nil {
		log.Fatalf("db init: %v", err)
	}

	srv := httptest.NewServer(api.NewRouter(conn))
	defer srv.Close()

	issueBody, _ := json.Marshal(map[string]string{"name": "smoketest"})
	resp, err := http.Post(srv.URL+"/api/v1/keys", "application/json", bytes.NewReader(issueBody))
	if err != nil {
		log.Fatal(err)
	}
	var issued api.IssueKeyResult
	if err := json.NewDecoder(resp.Body).Decode(&issued); err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	if issued.APIKey == "" {
		log.Fatal("FAIL: bootstrap key issuance returned empty key")
	}
	fmt.Println("PASS: bootstrap key issued")

	authed := func(method, path string, body []byte) *http.Response {
		req, _ := http.NewRequest(method, srv.URL+path, bytes.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+issued.APIKey)
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		return resp
	}

	// create suite
	suiteBody, _ := json.Marshal(api.CreateSuiteInput{Name: "greeting-bot", Description: "checks the greeting output"})
	resp = authed(http.MethodPost, "/api/v1/suites", suiteBody)
	var suite db.TestSuite
	if err := json.NewDecoder(resp.Body).Decode(&suite); err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	if suite.ID == "" {
		log.Fatal("FAIL: suite creation returned empty id")
	}
	fmt.Println("PASS: suite created")

	// create case A: string rules
	caseABody, _ := json.Marshal(api.CreateCaseInput{
		Name: "says hello",
		Rules: []db.Rule{
			{Type: "contains", Value: "Hello"},
			{Type: "min_length", Value: "5"},
		},
	})
	resp = authed(http.MethodPost, "/api/v1/suites/"+suite.ID+"/cases", caseABody)
	var caseA db.TestCase
	if err := json.NewDecoder(resp.Body).Decode(&caseA); err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()

	// create case B: JSON rules
	caseBBody, _ := json.Marshal(api.CreateCaseInput{
		Name: "returns valid status JSON",
		Rules: []db.Rule{
			{Type: "json_valid"},
			{Type: "json_key_exists", Value: "status"},
		},
	})
	resp = authed(http.MethodPost, "/api/v1/suites/"+suite.ID+"/cases", caseBBody)
	var caseB db.TestCase
	if err := json.NewDecoder(resp.Body).Decode(&caseB); err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	fmt.Println("PASS: two cases created (string rules + JSON rules)")

	// run 1: both cases pass
	run1Body, _ := json.Marshal(api.RunSuiteInput{
		Source: "smoketest",
		Cases: []api.RunCaseInput{
			{CaseID: caseA.ID, Output: "Hello there, friend"},
			{CaseID: caseB.ID, Output: `{"status":"ok"}`},
		},
	})
	resp = authed(http.MethodPost, "/api/v1/suites/"+suite.ID+"/run", run1Body)
	var run1 api.RunResult
	if err := json.NewDecoder(resp.Body).Decode(&run1); err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	if !run1.Run.Passed || run1.Run.Score != 1.0 {
		log.Fatalf("FAIL: expected run1 to fully pass with score=1.0, got passed=%v score=%v", run1.Run.Passed, run1.Run.Score)
	}
	fmt.Println("PASS: run1 — both cases pass, score=1.0")

	// run 2: case B fails (invalid JSON)
	run2Body, _ := json.Marshal(api.RunSuiteInput{
		Source: "smoketest",
		Cases: []api.RunCaseInput{
			{CaseID: caseA.ID, Output: "Hello there, friend"},
			{CaseID: caseB.ID, Output: "not json at all"},
		},
	})
	resp = authed(http.MethodPost, "/api/v1/suites/"+suite.ID+"/run", run2Body)
	var run2 api.RunResult
	if err := json.NewDecoder(resp.Body).Decode(&run2); err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	if run2.Run.Passed || run2.Run.Score != 0.5 {
		log.Fatalf("FAIL: expected run2 to fail with score=0.5, got passed=%v score=%v", run2.Run.Passed, run2.Run.Score)
	}
	fmt.Println("PASS: run2 — case B fails on invalid JSON, score=0.5")

	// run history should have 2 entries
	resp = authed(http.MethodGet, "/api/v1/suites/"+suite.ID+"/runs", nil)
	var runs []db.TestRun
	if err := json.NewDecoder(resp.Body).Decode(&runs); err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	if len(runs) != 2 {
		log.Fatalf("FAIL: expected 2 runs in history, got %d", len(runs))
	}
	fmt.Println("PASS: run history contains both runs")

	// run detail should include per-case rule breakdown
	resp = authed(http.MethodGet, "/api/v1/runs/"+run2.Run.ID, nil)
	var detail api.RunResult
	if err := json.NewDecoder(resp.Body).Decode(&detail); err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	if len(detail.Results) != 2 {
		log.Fatalf("FAIL: expected 2 case results in run detail, got %d", len(detail.Results))
	}
	fmt.Println("PASS: run detail returns per-case results")

	// unknown case_id must error, not silently skip
	badRunBody, _ := json.Marshal(api.RunSuiteInput{
		Source: "smoketest",
		Cases:  []api.RunCaseInput{{CaseID: "does-not-exist", Output: "x"}},
	})
	resp = authed(http.MethodPost, "/api/v1/suites/"+suite.ID+"/run", badRunBody)
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		log.Fatalf("FAIL: expected 400 for unknown case_id, got %d", resp.StatusCode)
	}
	fmt.Println("PASS: unknown case_id is rejected instead of silently skipped")

	fmt.Println("SMOKE TEST OK")
}
