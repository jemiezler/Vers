package llm

import (
	"strings"
	"testing"
)

func TestNewClient_DefaultsToStub(t *testing.T) {
	client, err := NewClient(Config{})
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}
	if _, ok := client.(*StubClient); !ok {
		t.Fatalf("client type = %T, want *StubClient", client)
	}
}

func TestNewClient_Stub(t *testing.T) {
	client, err := NewClient(Config{Provider: "StuB"})
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}
	if _, ok := client.(*StubClient); !ok {
		t.Fatalf("client type = %T, want *StubClient", client)
	}
}

func TestNewClient_Ollama(t *testing.T) {
	client, err := NewClient(Config{Provider: "ollama"})
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}
	if _, ok := client.(*OllamaClient); !ok {
		t.Fatalf("client type = %T, want *OllamaClient", client)
	}
}

func TestNewClient_UnsupportedProvider(t *testing.T) {
	_, err := NewClient(Config{Provider: "nope"})
	if err == nil {
		t.Fatal("NewClient returned nil error for unsupported provider")
	}
	if !strings.Contains(err.Error(), "unsupported") {
		t.Fatalf("error = %q, want contains 'unsupported'", err.Error())
	}
}
