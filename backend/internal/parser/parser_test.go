package parser

import "testing"

func TestParseGoMod(t *testing.T) {
	content := []byte(`module sample

go 1.22

require (
	github.com/gin-gonic/gin v1.10.0
	github.com/stretchr/testify v1.9.0 // indirect
)

require golang.org/x/sync v0.7.0
`)

	manifest, err := Parse("go.mod", content)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if manifest.Kind != "go.mod" {
		t.Fatalf("Kind = %q, want go.mod", manifest.Kind)
	}

	want := map[string]string{
		"github.com/gin-gonic/gin":    "v1.10.0",
		"github.com/stretchr/testify": "v1.9.0",
		"golang.org/x/sync":           "v0.7.0",
	}
	assertDependencies(t, manifest.Dependencies, want)
}

func TestParsePackageJSON(t *testing.T) {
	content := []byte(`{
  "dependencies": {
    "express": "^4.18.3"
  },
  "devDependencies": {
    "vitest": "^1.6.0"
  }
}`)

	manifest, err := Parse("package.json", content)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if manifest.Kind != "package.json" {
		t.Fatalf("Kind = %q, want package.json", manifest.Kind)
	}

	want := map[string]string{
		"express": "^4.18.3",
		"vitest":  "^1.6.0",
	}
	assertDependencies(t, manifest.Dependencies, want)
}

func TestParseUnsupportedManifest(t *testing.T) {
	_, err := Parse("requirements.txt", []byte("requests==2.32.0"))
	if err == nil {
		t.Fatal("Parse returned nil error for unsupported manifest")
	}
}

func assertDependencies(t *testing.T, got []Dependency, want map[string]string) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("dependency count = %d, want %d: %#v", len(got), len(want), got)
	}

	for _, dep := range got {
		version, ok := want[dep.Name]
		if !ok {
			t.Fatalf("unexpected dependency %q", dep.Name)
		}
		if dep.Version != version {
			t.Fatalf("dependency %q version = %q, want %q", dep.Name, dep.Version, version)
		}
		if dep.Source == "" {
			t.Fatalf("dependency %q has empty source", dep.Name)
		}
	}
}
