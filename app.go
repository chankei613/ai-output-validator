package main

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	"github.com/chankei613/ai-output-validator/internal/api"
	"github.com/chankei613/ai-output-validator/internal/db"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const apiAddr = "127.0.0.1:8426"

// App はWailsのバインディング。実処理は internal/api.Server が持っている。
type App struct {
	ctx    context.Context
	server *api.Server
	srv    *http.Server
	ready  bool
}

func NewApp() *App { return &App{} }

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	dataDir := appDataDir()
	if err := os.MkdirAll(dataDir, 0o700); err != nil {
		runtime.LogErrorf(ctx, "data dir error: %s", err)
		return
	}

	conn, err := db.Init(filepath.Join(dataDir, "ai-output-validator.db"))
	if err != nil {
		runtime.LogErrorf(ctx, "db init error: %s", err)
		return
	}
	a.server = api.New(conn)

	a.srv = &http.Server{Addr: apiAddr, Handler: a.server.Router()}
	go func() {
		if err := a.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			runtime.LogErrorf(ctx, "api server error: %s", err)
		}
	}()

	a.ready = true
	runtime.LogInfof(ctx, "AI Output Validator ready (api: http://%s, data: %s)", apiAddr, dataDir)
}

func (a *App) shutdown(ctx context.Context) {
	if a.srv != nil {
		_ = a.srv.Close()
	}
	if a.server != nil {
		if sqlDB, err := a.server.DB.DB(); err == nil {
			_ = sqlDB.Close()
		}
	}
}

var errNotReady = &notReadyError{}

type notReadyError struct{}

func (e *notReadyError) Error() string { return "app not ready — check startup logs" }

// ─── フロントエンドへ公開するメソッド ──────────────────────────────────────────

func (a *App) GetAppVersion() string {
	return AppVersion
}

func (a *App) GetAPIURL() string {
	return "http://" + apiAddr
}

func (a *App) ListSuites() ([]db.TestSuite, error) {
	if !a.ready {
		return nil, errNotReady
	}
	return a.server.ListSuites()
}

func (a *App) GetSuite(id string) (db.TestSuite, error) {
	if !a.ready {
		return db.TestSuite{}, errNotReady
	}
	return a.server.GetSuite(id)
}

func (a *App) CreateSuite(name, description string) (db.TestSuite, error) {
	if !a.ready {
		return db.TestSuite{}, errNotReady
	}
	return a.server.CreateSuite(api.CreateSuiteInput{Name: name, Description: description})
}

func (a *App) DeleteSuite(id string) error {
	if !a.ready {
		return errNotReady
	}
	return a.server.DeleteSuite(id)
}

func (a *App) CreateCase(suiteID, name string, rules []db.Rule) (db.TestCase, error) {
	if !a.ready {
		return db.TestCase{}, errNotReady
	}
	return a.server.CreateCase(suiteID, api.CreateCaseInput{Name: name, Rules: rules})
}

func (a *App) ListCases(suiteID string) ([]db.TestCase, error) {
	if !a.ready {
		return nil, errNotReady
	}
	return a.server.ListCases(suiteID)
}

func (a *App) DeleteCase(caseID string) error {
	if !a.ready {
		return errNotReady
	}
	return a.server.DeleteCase(caseID)
}

func (a *App) RunSuite(suiteID, source string, cases []api.RunCaseInput) (api.RunResult, error) {
	if !a.ready {
		return api.RunResult{}, errNotReady
	}
	return a.server.RunSuite(suiteID, api.RunSuiteInput{Source: source, Cases: cases})
}

func (a *App) ListRuns(suiteID string) ([]db.TestRun, error) {
	if !a.ready {
		return nil, errNotReady
	}
	return a.server.ListRuns(suiteID)
}

func (a *App) GetRun(id string) (api.RunResult, error) {
	if !a.ready {
		return api.RunResult{}, errNotReady
	}
	return a.server.GetRun(id)
}

func (a *App) ListKeys() ([]db.AgentKey, error) {
	if !a.ready {
		return nil, errNotReady
	}
	return a.server.ListKeys()
}

func (a *App) IssueKey(name string) (api.IssueKeyResult, error) {
	if !a.ready {
		return api.IssueKeyResult{}, errNotReady
	}
	return a.server.IssueKey(name)
}

func (a *App) RevokeKey(id string) error {
	if !a.ready {
		return errNotReady
	}
	return a.server.RevokeKey(id)
}

// Quit はアプリを完全終了する（Settings 画面から呼ぶ）。
func (a *App) Quit() {
	runtime.Quit(a.ctx)
}

func appDataDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return filepath.Join(home, ".ai-output-validator")
}
