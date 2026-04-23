package scan

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jemiezler/Vers/backend/internal/review"
)

type FileReport struct {
	Path       string         `json:"path"`
	Filename   string         `json:"filename"`
	DurationMs int64          `json:"durationMs"`
	Result     *review.Result `json:"result,omitempty"`
	Error      string         `json:"error,omitempty"`
}

type Report struct {
	Root       string       `json:"root"`
	StartedAt  time.Time    `json:"startedAt"`
	FinishedAt time.Time    `json:"finishedAt"`
	FileCount  int          `json:"fileCount"`
	ErrorCount int          `json:"errorCount"`
	Files      []FileReport `json:"files"`
}

func ScanPath(root string, maxFiles int, service *review.Service) (Report, error) {
	if maxFiles <= 0 {
		return Report{}, errors.New("maxFiles must be > 0")
	}
	start := time.Now()

	abs, err := filepath.Abs(root)
	if err != nil {
		return Report{}, err
	}

	files := make([]FileReport, 0)
	errCount := 0

	walkErr := filepath.WalkDir(abs, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			errCount++
			files = append(files, FileReport{
				Path:  path,
				Error: err.Error(),
			})
			return nil
		}

		if d.IsDir() {
			base := strings.ToLower(d.Name())
			if base == ".git" || base == "node_modules" || base == "vendor" || base == "dist" || base == "build" {
				return filepath.SkipDir
			}
			return nil
		}

		filename := filepath.Base(path)
		if !isSupportedManifest(filename) {
			return nil
		}

		if len(files) >= maxFiles {
			return fs.SkipAll
		}

		report := scanOne(path, filename, service)
		if report.Error != "" {
			errCount++
		}
		files = append(files, report)
		return nil
	})
	if walkErr != nil {
		return Report{}, walkErr
	}

	finish := time.Now()
	return Report{
		Root:       abs,
		StartedAt:  start,
		FinishedAt: finish,
		FileCount:  len(files),
		ErrorCount: errCount,
		Files:      files,
	}, nil
}

func isSupportedManifest(filename string) bool {
	switch strings.ToLower(filename) {
	case "go.mod", "package.json":
		return true
	default:
		return false
	}
}

func scanOne(path string, filename string, service *review.Service) FileReport {
	start := time.Now()
	content, err := os.ReadFile(path)
	if err != nil {
		return FileReport{
			Path:       path,
			Filename:   filename,
			DurationMs: sinceMs(start),
			Error:      err.Error(),
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	result, err := service.Run(ctx, review.Request{
		Filename: filename,
		Content:  string(content),
	})
	if err != nil {
		return FileReport{
			Path:       path,
			Filename:   filename,
			DurationMs: sinceMs(start),
			Error:      err.Error(),
		}
	}

	return FileReport{
		Path:       path,
		Filename:   filename,
		DurationMs: sinceMs(start),
		Result:     &result,
	}
}

func sinceMs(start time.Time) int64 {
	return time.Since(start).Milliseconds()
}
