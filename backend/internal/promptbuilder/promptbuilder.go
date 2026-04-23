package promptbuilder

import (
	"fmt"
	"strings"

	"vers/backend/internal/contextbuilder"
)

func Build(ctx contextbuilder.ReviewContext) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Review a project using %s dependencies.\n", ctx.ManifestKind)
	b.WriteString("Focus on version-specific risks, breaking changes, and insecure usage.\n")
	b.WriteString("Output MUST be plain ASCII only (no emojis, no Unicode bullets, no special symbols).\n")
	b.WriteString("Every factual claim MUST include a citation to a provided source URL in the form: [source: <url>].\n")
	b.WriteString("If the provided docs do not support a claim, explicitly say: unknown from provided docs.\n\n")

	for _, dep := range ctx.Dependencies {
		fmt.Fprintf(&b, "- %s %s\n", dep.Dependency.Name, dep.Dependency.Version)
		for _, doc := range dep.Docs {
			fmt.Fprintf(&b, "  doc_source: %s\n", asciiOnly(doc.Source))
			fmt.Fprintf(&b, "  doc_text: %s\n", asciiOnly(doc.Text))
		}
	}

	return b.String()
}

func asciiOnly(s string) string {
	// Keep only printable ASCII plus whitespace, to prevent the model from echoing
	// non-ASCII characters that render poorly in terminals.
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '\n' || c == '\r' || c == '\t':
			b.WriteByte(c)
		case c >= 0x20 && c <= 0x7E:
			b.WriteByte(c)
		default:
			// drop
		}
	}

	out := strings.TrimSpace(b.String())
	const max = 4000
	if len(out) > max {
		out = out[:max]
	}
	return out
}
