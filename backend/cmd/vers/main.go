package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/jemiezler/Vers/backend/internal/config"
	"github.com/jemiezler/Vers/backend/internal/review"
	"github.com/jemiezler/Vers/backend/internal/scan"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "scan":
		scanCmd(os.Args[2:])
	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage:")
	fmt.Fprintln(os.Stderr, "  vers scan <path> [--output <file>] [--max-files N]")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "env defaults (override as needed):")
	fmt.Fprintln(os.Stderr, "  VERS_LLM_PROVIDER=stub|ollama")
	fmt.Fprintln(os.Stderr, "  VERS_DOCS_PROVIDER=stub|pkg_go_dev")
}

func scanCmd(args []string) {
	fs := flag.NewFlagSet("scan", flag.ExitOnError)
	output := fs.String("output", "", "write JSON report to file (default stdout)")
	maxFiles := fs.Int("max-files", 200, "maximum manifests to scan")
	failOnError := fs.Bool("fail-on-error", true, "return non-zero exit code if any file scan errored")
	_ = fs.Parse(args)

	if fs.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "missing <path>")
		fs.Usage()
		os.Exit(2)
	}
	root := fs.Arg(0)

	cfg := config.Load()
	service, err := review.NewServiceFromConfig(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "configure review service: %v\n", err)
		os.Exit(1)
	}

	report, err := scan.ScanPath(root, *maxFiles, service)
	if err != nil {
		fmt.Fprintf(os.Stderr, "scan failed: %v\n", err)
		os.Exit(1)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if *output != "" {
		f, err := os.Create(*output)
		if err != nil {
			fmt.Fprintf(os.Stderr, "create output file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		enc = json.NewEncoder(f)
		enc.SetIndent("", "  ")
	}

	if err := enc.Encode(report); err != nil {
		fmt.Fprintf(os.Stderr, "write report: %v\n", err)
		os.Exit(1)
	}

	if *failOnError && report.ErrorCount > 0 {
		os.Exit(3)
	}
}
