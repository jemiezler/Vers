package promptbuilder

import (
	"fmt"
	"strings"

	"vers/backend/internal/contextbuilder"
)

func Build(ctx contextbuilder.ReviewContext) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Review a project using %s dependencies.\n", ctx.ManifestKind)
	b.WriteString("Focus on version-specific risks, breaking changes, and insecure usage.\n\n")

	for _, dep := range ctx.Dependencies {
		fmt.Fprintf(&b, "- %s %s\n", dep.Dependency.Name, dep.Dependency.Version)
		for _, doc := range dep.Docs {
			fmt.Fprintf(&b, "  doc: %s\n", doc)
		}
	}

	return b.String()
}
