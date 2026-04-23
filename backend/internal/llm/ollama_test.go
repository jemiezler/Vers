package llm

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestOllamaClientComplete_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/api/generate" {
			t.Fatalf("path = %s, want /api/generate", r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); !strings.Contains(got, "application/json") {
			t.Fatalf("Content-Type = %q, want application/json", got)
		}

		raw, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("ReadAll returned error: %v", err)
		}

		var req ollamaGenerateRequest
		if err := json.Unmarshal(raw, &req); err != nil {
			t.Fatalf("Unmarshal returned error: %v", err)
		}
		if req.Model != "gemma3" {
			t.Fatalf("model = %q, want gemma3", req.Model)
		}
		if req.Prompt != "hi" {
			t.Fatalf("prompt = %q, want hi", req.Prompt)
		}
		if req.Stream {
			t.Fatalf("stream = %v, want false", req.Stream)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(ollamaGenerateResponse{Response: "ok"})
	}))
	defer server.Close()

	client := NewOllamaClient(server.URL, "gemma3")
	client.httpClient.Timeout = 2 * time.Second

	out, err := client.Complete(context.Background(), "hi")
	if err != nil {
		t.Fatalf("Complete returned error: %v", err)
	}
	if out != "ok" {
		t.Fatalf("out = %q, want ok", out)
	}
}

func TestOllamaClientComplete_HTTPError_JSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ollamaGenerateResponse{Error: "bad model"})
	}))
	defer server.Close()

	client := NewOllamaClient(server.URL, "gemma3")
	client.httpClient.Timeout = 2 * time.Second

	_, err := client.Complete(context.Background(), "hi")
	if err == nil {
		t.Fatal("Complete returned nil error")
	}
	if !strings.Contains(err.Error(), "400") || !strings.Contains(err.Error(), "bad model") {
		t.Fatalf("error = %q, want contains 400 and bad model", err.Error())
	}
}

func TestOllamaClientComplete_HTTPError_PlainText(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("boom"))
	}))
	defer server.Close()

	client := NewOllamaClient(server.URL, "gemma3")
	client.httpClient.Timeout = 2 * time.Second

	_, err := client.Complete(context.Background(), "hi")
	if err == nil {
		t.Fatal("Complete returned nil error")
	}
	if !strings.Contains(err.Error(), "500") || !strings.Contains(err.Error(), "boom") {
		t.Fatalf("error = %q, want contains 500 and boom", err.Error())
	}
}
