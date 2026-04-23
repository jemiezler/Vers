package vectordb

import "vers/backend/internal/ingestion/embedder"

type Chunk struct {
	Library string             `json:"library"`
	Version string             `json:"version"`
	Text    string             `json:"text"`
	Vector  embedder.Embedding `json:"-"`
}

type MemoryStore struct {
	chunks []Chunk
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{chunks: make([]Chunk, 0)}
}

func (s *MemoryStore) Add(chunk Chunk) {
	s.chunks = append(s.chunks, chunk)
}

func (s *MemoryStore) Search(library string) []Chunk {
	results := make([]Chunk, 0)
	for _, chunk := range s.chunks {
		if chunk.Library == library {
			results = append(results, chunk)
		}
	}

	return results
}
