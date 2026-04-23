package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"vers/backend/internal/review"
)

func TestHealth(t *testing.T) {
	handler := NewHandler(review.NewService())
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)

	handler.Health(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if !strings.Contains(recorder.Body.String(), `"status":"ok"`) {
		t.Fatalf("body = %s, want health payload", recorder.Body.String())
	}
}

func TestCreateReview(t *testing.T) {
	handler := NewHandler(review.NewService())
	body := review.Request{
		Filename: "go.mod",
		Content:  "require github.com/gin-gonic/gin v1.10.0",
	}
	payload, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/reviews", bytes.NewReader(payload))

	handler.CreateReview(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body: %s", recorder.Code, http.StatusOK, recorder.Body.String())
	}

	var result review.Result
	if err := json.NewDecoder(recorder.Body).Decode(&result); err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}

	if len(result.Manifest.Dependencies) != 1 {
		t.Fatalf("dependency count = %d, want 1", len(result.Manifest.Dependencies))
	}
	if result.Manifest.Dependencies[0].Name != "github.com/gin-gonic/gin" {
		t.Fatalf("dependency name = %q, want github.com/gin-gonic/gin", result.Manifest.Dependencies[0].Name)
	}
	if len(result.Context.Dependencies) != 1 {
		t.Fatalf("context dependency count = %d, want 1", len(result.Context.Dependencies))
	}
	if len(result.Context.Dependencies[0].Docs) != 1 {
		t.Fatalf("doc chunk count = %d, want 1", len(result.Context.Dependencies[0].Docs))
	}
	if result.Context.Dependencies[0].Docs[0].Source == "" {
		t.Fatal("doc chunk source is empty")
	}
	if result.Context.Dependencies[0].Docs[0].Text == "" {
		t.Fatal("doc chunk text is empty")
	}
	if result.Prompt == "" {
		t.Fatal("prompt is empty")
	}
	if result.Review == "" {
		t.Fatal("review is empty")
	}
	if result.LLM == "" {
		t.Fatal("llmProvider is empty")
	}
}

func TestCreateReviewRejectsInvalidJSON(t *testing.T) {
	handler := NewHandler(review.NewService())
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/reviews", strings.NewReader("{"))

	handler.CreateReview(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
}
