package contextbuilder

import "vers/backend/internal/parser"

type DocChunk struct {
	Source string `json:"source"`
	Text   string `json:"text"`
}

type DependencyContext struct {
	Dependency parser.Dependency `json:"dependency"`
	Docs       []DocChunk        `json:"docs"`
}

type ReviewContext struct {
	ManifestKind string              `json:"manifestKind"`
	Dependencies []DependencyContext `json:"dependencies"`
}

func Build(manifest parser.Manifest, docs map[string][]DocChunk) ReviewContext {
	deps := make([]DependencyContext, 0, len(manifest.Dependencies))
	for _, dep := range manifest.Dependencies {
		deps = append(deps, DependencyContext{
			Dependency: dep,
			Docs:       docs[dep.Name],
		})
	}

	return ReviewContext{
		ManifestKind: manifest.Kind,
		Dependencies: deps,
	}
}
