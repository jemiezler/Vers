package review

import (
	"context"
	"errors"

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
}

type Service struct {
	fetcher  *fetcher.Fetcher
	embedder *embedder.Embedder
	store    *vectordb.MemoryStore
	llm      llm.Client
}

func NewService() *Service {
	return &Service{
		fetcher:  fetcher.New(),
		embedder: embedder.New(),
		store:    vectordb.NewMemoryStore(),
		llm:      llm.NewStubClient(),
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

	docs := make(map[string][]string, len(manifest.Dependencies))
	for _, dep := range manifest.Dependencies {
		doc := s.fetcher.Fetch(dep)
		relevant := scraper.ExtractRelevant(doc)
		markdown := converter.ToMarkdown(relevant)
		vector := s.embedder.Embed(markdown)
		s.store.Add(vectordb.Chunk{
			Library: dep.Name,
			Version: dep.Version,
			Text:    markdown,
			Vector:  vector,
		})

		for _, chunk := range s.store.Search(dep.Name) {
			docs[dep.Name] = append(docs[dep.Name], chunk.Text)
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
	}, nil
}
