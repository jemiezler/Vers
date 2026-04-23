package llm

import (
	"context"
	"strconv"
)

type Client interface {
	Complete(ctx context.Context, prompt string) (string, error)
}

type StubClient struct{}

func NewStubClient() *StubClient {
	return &StubClient{}
}

func (c *StubClient) Complete(_ context.Context, prompt string) (string, error) {
	return "Stub review generated from prompt length " + strconv.Itoa(len(prompt)) + ". Wire this to Ollama or vLLM next.", nil
}
