package promptbuilder

import (
	"testing"

	"github.com/jemiezler/Vers/backend/internal/contextbuilder"
	"github.com/jemiezler/Vers/backend/internal/parser"
)

func TestBuild_AsciiOnly(t *testing.T) {
	ctx := contextbuilder.ReviewContext{
		ManifestKind: "go.mod",
		Dependencies: []contextbuilder.DependencyContext{
			{
				Dependency: parser.Dependency{
					Name:    "example.com/lib",
					Version: "v1.2.3",
					Source:  "go.mod",
				},
				Docs: []contextbuilder.DocChunk{
					{Source: "https://example.com/docs", Text: "hello 😅 — café"},
				},
			},
		},
	}

	out := Build(ctx)
	if hasNonASCII(out) {
		t.Fatalf("prompt contains non-ASCII: %q", out)
	}
}

func hasNonASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= 0x80 {
			return true
		}
	}
	return false
}
