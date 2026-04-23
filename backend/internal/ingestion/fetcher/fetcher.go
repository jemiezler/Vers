package fetcher

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"vers/backend/internal/parser"
)

type Document struct {
	Dependency parser.Dependency
	URL        string
	Content    string
}

type Fetcher struct {
	docsProvider string
	pkgGoDevURL  string
	httpClient   *http.Client
}

type Config struct {
	DocsProvider string
	PkgGoDevURL  string
}

func New(cfg Config) *Fetcher {
	pkgGoDevURL := cfg.PkgGoDevURL
	if pkgGoDevURL == "" {
		pkgGoDevURL = "https://pkg.go.dev"
	}

	return &Fetcher{
		docsProvider: strings.ToLower(strings.TrimSpace(cfg.DocsProvider)),
		pkgGoDevURL:  strings.TrimRight(pkgGoDevURL, "/"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (f *Fetcher) Fetch(dep parser.Dependency) Document {
	switch f.docsProvider {
	case "", "stub":
		return stubDoc(dep)
	case "pkg_go_dev", "pkg.go.dev", "pkg":
		if dep.Source != "go.mod" {
			return stubDoc(dep)
		}
		return f.fetchPkgGoDev(dep)
	default:
		return stubDoc(dep)
	}
}

func stubDoc(dep parser.Dependency) Document {
	return Document{
		Dependency: dep,
		URL:        "local://docs/" + dep.Name,
		Content:    "Documentation placeholder for " + dep.Name + " " + dep.Version,
	}
}

func (f *Fetcher) fetchPkgGoDev(dep parser.Dependency) Document {
	url := f.pkgGoDevURL + "/" + dep.Name + "@" + dep.Version + "?tab=doc"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return stubDoc(dep)
	}
	req.Header.Set("Accept", "text/html")

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return stubDoc(dep)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return stubDoc(dep)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return Document{
			Dependency: dep,
			URL:        url,
			Content:    fmt.Sprintf("pkg.go.dev returned %d", resp.StatusCode),
		}
	}

	return Document{
		Dependency: dep,
		URL:        url,
		Content:    string(body),
	}
}
