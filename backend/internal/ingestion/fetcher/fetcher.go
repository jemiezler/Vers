package fetcher

import "vers/backend/internal/parser"

type Document struct {
	Dependency parser.Dependency
	URL        string
	Content    string
}

type Fetcher struct{}

func New() *Fetcher {
	return &Fetcher{}
}

func (f *Fetcher) Fetch(dep parser.Dependency) Document {
	return Document{
		Dependency: dep,
		URL:        "local://docs/" + dep.Name,
		Content:    "Documentation placeholder for " + dep.Name + " " + dep.Version,
	}
}
