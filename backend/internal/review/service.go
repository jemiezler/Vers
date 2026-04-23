package review

import (
	"context"
	"errors"

	"vers/backend/internal/config"
	"vers/backend/internal/contextbuilder"
	"vers/backend/internal/ingestion/converter"
	"vers/backend/internal/ingestion/embedder"
	"vers/backend/internal/ingestion/fetcher"
	"vers/backend/internal/ingestion/scraper"
	"vers/backend/internal/llm"
	"vers/backend/internal/parser"
	"vers/backend/internal/promptbuilder"
	"vers/backend/internal/vectordb"
)

type Request struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
}

type Result struct {
	Manifest parser.Manifest              `json:"manifest"`
	Context  contextbuilder.ReviewContext `json:"context"`
	Prompt   string                       `json:"prompt"`
	Review   string                       `json:"review"`
	LLM      string                       `json:"llmProvider"`
}

type Service struct {
	fetcher  *fetcher.Fetcher
	embedder *embedder.Embedder
	llm      llm.Client
	llmName  string
}

func NewService() *Service {
	service, err := NewServiceFromConfig(config.Load())
	if err != nil {
		return NewServiceWithLLM(
			"stub",
			llm.NewStubClient(),
			fetcher.New(fetcher.Config{DocsProvider: "stub"}),
		)
	}

	return service
}

func NewServiceFromConfig(cfg config.Config) (*Service, error) {
	client, err := llm.NewClient(llm.Config{
		Provider:    cfg.LLMProvider,
		OllamaURL:   cfg.OllamaURL,
		OllamaModel: cfg.OllamaModel,
	})
	if err != nil {
		return nil, err
	}

	name := cfg.LLMProvider
	if name == "" {
		name = "stub"
	}

	docsFetcher := fetcher.New(fetcher.Config{
		DocsProvider: cfg.DocsProvider,
		PkgGoDevURL:  cfg.PkgGoDevURL,
	})

	return NewServiceWithLLM(name, client, docsFetcher), nil
}

func NewServiceWithLLM(name string, client llm.Client, docsFetcher *fetcher.Fetcher) *Service {
	return &Service{
		fetcher:  docsFetcher,
		embedder: embedder.New(),
		llm:      client,
		llmName:  name,
	}
}

func (s *Service) Run(ctx context.Context, req Request) (Result, error) {
	if req.Filename == "" || req.Content == "" {
		return Result{}, errors.New("filename and content are required")
	}

	manifest, err := parser.Parse(req.Filename, []byte(req.Content))
	if err != nil {
		return Result{}, err
	}

	// Keep the scaffold request-scoped so repeated calls don't duplicate docs.
	store := vectordb.NewMemoryStore()

	docs := make(map[string][]contextbuilder.DocChunk, len(manifest.Dependencies))
	for _, dep := range manifest.Dependencies {
		doc := s.fetcher.Fetch(dep)
		relevant := scraper.ExtractRelevant(doc)
		markdown := converter.ToMarkdown(relevant)
		vector := s.embedder.Embed(markdown)
		store.Add(vectordb.Chunk{
			Library: dep.Name,
			Version: dep.Version,
			Source:  doc.URL,
			Text:    markdown,
			Vector:  vector,
		})

		for _, chunk := range store.Search(dep.Name) {
			docs[dep.Name] = append(docs[dep.Name], contextbuilder.DocChunk{
				Source: chunk.Source,
				Text:   chunk.Text,
			})
		}
	}

	reviewContext := contextbuilder.Build(manifest, docs)
	prompt := promptbuilder.Build(reviewContext)
	response, err := s.llm.Complete(ctx, prompt)
	if err != nil {
		return Result{}, err
	}

	return Result{
		Manifest: manifest,
		Context:  reviewContext,
		Prompt:   prompt,
		Review:   response,
		LLM:      s.llmName,
	}, nil
}
