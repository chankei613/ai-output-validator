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
	errSuiteNameRequired = &apiError{"name is required"}
	errSuiteNotFound     = &apiError{"suite not found"}
	errCaseNameRequired  = &apiError{"name is required"}
	errCaseNotFound      = &apiError{"case not found"}
)

type CreateSuiteInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s *Server) CreateSuite(in CreateSuiteInput) (db.TestSuite, error) {
	if in.Name == "" {
		return db.TestSuite{}, errSuiteNameRequired
	}
	now := time.Now()
	suite := db.TestSuite{
		ID:          uuid.NewString(),
		Name:        in.Name,
		Description: in.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.DB.Create(&suite).Error; err != nil {
		return db.TestSuite{}, err
	}
	return suite, nil
}

func (s *Server) ListSuites() ([]db.TestSuite, error) {
	var suites []db.TestSuite
	err := s.DB.Order("created_at asc").Find(&suites).Error
	return suites, err
}

func (s *Server) GetSuite(id string) (db.TestSuite, error) {
	var suite db.TestSuite
	if err := s.DB.First(&suite, "id = ?", id).Error; err != nil {
		return db.TestSuite{}, errSuiteNotFound
	}
	return suite, nil
}

func (s *Server) DeleteSuite(id string) error {
	res := s.DB.Delete(&db.TestSuite{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errSuiteNotFound
	}
	return s.DB.Where("suite_id = ?", id).Delete(&db.TestCase{}).Error
}

type CreateCaseInput struct {
	Name  string    `json:"name"`
	Rules []db.Rule `json:"rules"`
}

func (s *Server) CreateCase(suiteID string, in CreateCaseInput) (db.TestCase, error) {
	if in.Name == "" {
		return db.TestCase{}, errCaseNameRequired
	}
	if _, err := s.GetSuite(suiteID); err != nil {
		return db.TestCase{}, err
	}
	tc := db.TestCase{
		ID:        uuid.NewString(),
		SuiteID:   suiteID,
		Name:      in.Name,
		Rules:     in.Rules,
		CreatedAt: time.Now(),
	}
	if err := s.DB.Create(&tc).Error; err != nil {
		return db.TestCase{}, err
	}
	return tc, nil
}

func (s *Server) ListCases(suiteID string) ([]db.TestCase, error) {
	var cases []db.TestCase
	err := s.DB.Where("suite_id = ?", suiteID).Order("created_at asc").Find(&cases).Error
	return cases, err
}

func (s *Server) DeleteCase(id string) error {
	res := s.DB.Delete(&db.TestCase{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errCaseNotFound
	}
	return nil
}

// ─── HTTP handlers ──────────────────────────────────────────────────────

func (s *Server) httpCreateSuite(w http.ResponseWriter, r *http.Request) {
	var in CreateSuiteInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	suite, err := s.CreateSuite(in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusCreated, suite)
}

func (s *Server) httpListSuites(w http.ResponseWriter, r *http.Request) {
	suites, err := s.ListSuites()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, suites)
}

func (s *Server) httpGetSuite(w http.ResponseWriter, r *http.Request) {
	suite, err := s.GetSuite(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, suite)
}

func (s *Server) httpDeleteSuite(w http.ResponseWriter, r *http.Request) {
	if err := s.DeleteSuite(chi.URLParam(r, "id")); err != nil {
		if err == errSuiteNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) httpCreateCase(w http.ResponseWriter, r *http.Request) {
	var in CreateCaseInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	tc, err := s.CreateCase(chi.URLParam(r, "id"), in)
	if err != nil {
		status := http.StatusBadRequest
		if err == errSuiteNotFound {
			status = http.StatusNotFound
		}
		http.Error(w, err.Error(), status)
		return
	}
	writeJSON(w, http.StatusCreated, tc)
}

func (s *Server) httpListCases(w http.ResponseWriter, r *http.Request) {
	cases, err := s.ListCases(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, cases)
}

func (s *Server) httpDeleteCase(w http.ResponseWriter, r *http.Request) {
	if err := s.DeleteCase(chi.URLParam(r, "caseId")); err != nil {
		if err == errCaseNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
