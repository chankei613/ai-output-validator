package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/chankei613/ai-output-validator/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

var (
	errRunNotFound     = &apiError{"run not found"}
	errNoCasesProvided = &apiError{"cases must not be empty"}
)

type RunCaseInput struct {
	CaseID string `json:"case_id"`
	Output string `json:"output"`
}

type RunSuiteInput struct {
	Source string         `json:"source"`
	Cases  []RunCaseInput `json:"cases"`
}

type RunResult struct {
	Run     db.TestRun      `json:"run"`
	Results []db.CaseResult `json:"results"`
}

// RunSuite はSuiteに属するTestCaseそれぞれについて、外部から届いた出力(Output)を
// 決定論的ルールで検証し、TestRun+CaseResultとして永続化する。
// 未知のCaseIDが渡された場合はスキップせずエラーにする（テストの取りこぼしを防ぐため）。
func (s *Server) RunSuite(suiteID string, in RunSuiteInput) (RunResult, error) {
	if _, err := s.GetSuite(suiteID); err != nil {
		return RunResult{}, err
	}
	if len(in.Cases) == 0 {
		return RunResult{}, errNoCasesProvided
	}

	cases, err := s.ListCases(suiteID)
	if err != nil {
		return RunResult{}, err
	}
	caseByID := make(map[string]db.TestCase, len(cases))
	for _, c := range cases {
		caseByID[c.ID] = c
	}

	startedAt := time.Now()
	runID := uuid.NewString()

	results := make([]db.CaseResult, 0, len(in.Cases))
	passedCount := 0

	for _, rc := range in.Cases {
		tc, ok := caseByID[rc.CaseID]
		if !ok {
			return RunResult{}, &apiError{"unknown case_id: " + rc.CaseID}
		}
		passed, ruleResults := evaluateCase(tc.Rules, rc.Output)
		if passed {
			passedCount++
		}
		results = append(results, db.CaseResult{
			ID:          uuid.NewString(),
			TestRunID:   runID,
			TestCaseID:  tc.ID,
			CaseName:    tc.Name,
			Output:      rc.Output,
			Passed:      passed,
			RuleResults: ruleResults,
		})
	}

	score := float64(passedCount) / float64(len(results))
	run := db.TestRun{
		ID:         runID,
		SuiteID:    suiteID,
		Source:     in.Source,
		Passed:     passedCount == len(results),
		Score:      score,
		StartedAt:  startedAt,
		FinishedAt: time.Now(),
	}

	if err := s.DB.Create(&run).Error; err != nil {
		return RunResult{}, err
	}
	if err := s.DB.Create(&results).Error; err != nil {
		return RunResult{}, err
	}

	return RunResult{Run: run, Results: results}, nil
}

func (s *Server) ListRuns(suiteID string) ([]db.TestRun, error) {
	var runs []db.TestRun
	err := s.DB.Where("suite_id = ?", suiteID).Order("started_at asc").Find(&runs).Error
	return runs, err
}

func (s *Server) GetRun(id string) (RunResult, error) {
	var run db.TestRun
	if err := s.DB.First(&run, "id = ?", id).Error; err != nil {
		return RunResult{}, errRunNotFound
	}
	var results []db.CaseResult
	if err := s.DB.Where("test_run_id = ?", id).Find(&results).Error; err != nil {
		return RunResult{}, err
	}
	return RunResult{Run: run, Results: results}, nil
}

// ─── HTTP handlers ──────────────────────────────────────────────────────

func (s *Server) httpRunSuite(w http.ResponseWriter, r *http.Request) {
	var in RunSuiteInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	result, err := s.RunSuite(chi.URLParam(r, "id"), in)
	if err != nil {
		status := http.StatusBadRequest
		if err == errSuiteNotFound {
			status = http.StatusNotFound
		}
		http.Error(w, err.Error(), status)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (s *Server) httpListRuns(w http.ResponseWriter, r *http.Request) {
	runs, err := s.ListRuns(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, runs)
}

func (s *Server) httpGetRun(w http.ResponseWriter, r *http.Request) {
	result, err := s.GetRun(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, result)
}
