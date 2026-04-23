package llm

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

type Client interface {
	Complete(ctx context.Context, prompt string) (string, error)
}

type Config struct {
	Provider    string
	OllamaURL   string
	OllamaModel string
}

func NewClient(cfg Config) (Client, error) {
	switch strings.ToLower(cfg.Provider) {
	case "", "stub":
		return NewStubClient(), nil
	case "ollama":
		return NewOllamaClient(cfg.OllamaURL, cfg.OllamaModel), nil
	default:
		return nil, fmt.Errorf("unsupported llm provider %q", cfg.Provider)
	}
}

type StubClient struct{}

func NewStubClient() *StubClient {
	return &StubClient{}
}

func (c *StubClient) Complete(_ context.Context, prompt string) (string, error) {
	return "Stub review generated from prompt length " + strconv.Itoa(len(prompt)) + ". Wire this to Ollama or vLLM next.", nil
}
