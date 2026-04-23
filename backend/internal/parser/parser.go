package parser

import (
	"encoding/json"
	"errors"
	"strings"
)

type Dependency struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Source  string `json:"source"`
}

type Manifest struct {
	Kind         string       `json:"kind"`
	Dependencies []Dependency `json:"dependencies"`
}

func Parse(filename string, content []byte) (Manifest, error) {
	switch {
	case strings.HasSuffix(filename, "go.mod"):
		return parseGoMod(content), nil
	case strings.HasSuffix(filename, "package.json"):
		return parsePackageJSON(content)
	default:
		return Manifest{}, errors.New("unsupported manifest type")
	}
}

func parseGoMod(content []byte) Manifest {
	lines := strings.Split(string(content), "\n")
	deps := make([]Dependency, 0)
	inBlock := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		if strings.HasPrefix(line, "require (") {
			inBlock = true
			continue
		}
		if inBlock && line == ")" {
			inBlock = false
			continue
		}
		if strings.HasPrefix(line, "require ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "require "))
		} else if !inBlock {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 2 {
			deps = append(deps, Dependency{Name: fields[0], Version: fields[1], Source: "go.mod"})
		}
	}

	return Manifest{Kind: "go.mod", Dependencies: deps}
}

func parsePackageJSON(content []byte) (Manifest, error) {
	var raw struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}
	if err := json.Unmarshal(content, &raw); err != nil {
		return Manifest{}, err
	}

	deps := make([]Dependency, 0, len(raw.Dependencies)+len(raw.DevDependencies))
	for name, version := range raw.Dependencies {
		deps = append(deps, Dependency{Name: name, Version: version, Source: "package.json"})
	}
	for name, version := range raw.DevDependencies {
		deps = append(deps, Dependency{Name: name, Version: version, Source: "package.json:dev"})
	}

	return Manifest{Kind: "package.json", Dependencies: deps}, nil
}
