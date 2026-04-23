package embedder

type Embedding []float32

type Embedder struct{}

func New() *Embedder {
	return &Embedder{}
}

func (e *Embedder) Embed(text string) Embedding {
	return Embedding{float32(len(text))}
}
