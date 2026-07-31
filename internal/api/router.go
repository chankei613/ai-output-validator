package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// Server は全ロジックの実体。HTTPハンドラーとWailsネイティブバインディングの
// 両方がこの同じ Server のメソッドを呼ぶことで、UIとAPIの挙動がズレないようにする。
type Server struct {
	DB *gorm.DB
}

func New(conn *gorm.DB) *Server {
	return &Server{DB: conn}
}

func (s *Server) Router() http.Handler {
	r := chi.NewRouter()

	r.Route("/api/v1/keys", func(r chi.Router) {
		r.Use(APIKeyAuth(s.DB, "/api/v1/keys"))
		r.Post("/", s.httpIssueKey)
		r.Get("/", s.httpListKeys)
		r.Delete("/{id}", s.httpRevokeKey)
	})

	r.Route("/api/v1/suites", func(r chi.Router) {
		r.Use(APIKeyAuth(s.DB))
		r.Post("/", s.httpCreateSuite)
		r.Get("/", s.httpListSuites)
		r.Get("/{id}", s.httpGetSuite)
		r.Delete("/{id}", s.httpDeleteSuite)
		r.Post("/{id}/cases", s.httpCreateCase)
		r.Get("/{id}/cases", s.httpListCases)
		r.Delete("/{id}/cases/{caseId}", s.httpDeleteCase)
		r.Post("/{id}/run", s.httpRunSuite)
		r.Get("/{id}/runs", s.httpListRuns)
	})

	r.Route("/api/v1/runs", func(r chi.Router) {
		r.Use(APIKeyAuth(s.DB))
		r.Get("/{id}", s.httpGetRun)
	})

	return r
}

// NewRouter はcmd/ovserve（単体HTTPサーバー）向けの簡易コンストラクタ。
func NewRouter(conn *gorm.DB) http.Handler {
	return New(conn).Router()
}
