package config

import "os"

type Config struct {
	HTTPAddr     string
	LLMProvider  string
	OllamaURL    string
	OllamaModel  string
	DocsProvider string
	PkgGoDevURL  string
}

func Load() Config {
	addr := os.Getenv("VERS_HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	provider := os.Getenv("VERS_LLM_PROVIDER")
	if provider == "" {
		provider = "stub"
	}

	ollamaURL := os.Getenv("VERS_OLLAMA_URL")
	if ollamaURL == "" {
		ollamaURL = "http://localhost:11434"
	}

	ollamaModel := os.Getenv("VERS_OLLAMA_MODEL")
	if ollamaModel == "" {
		ollamaModel = "gemma3"
	}

	docsProvider := os.Getenv("VERS_DOCS_PROVIDER")
	if docsProvider == "" {
		docsProvider = "stub"
	}

	pkgGoDevURL := os.Getenv("VERS_PKG_GO_DEV_URL")
	if pkgGoDevURL == "" {
		pkgGoDevURL = "https://pkg.go.dev"
	}

	return Config{
		HTTPAddr:     addr,
		LLMProvider:  provider,
		OllamaURL:    ollamaURL,
		OllamaModel:  ollamaModel,
		DocsProvider: docsProvider,
		PkgGoDevURL:  pkgGoDevURL,
	}
}
