# Architecture

Vers is organized around the pipeline in `README.md`.

1. The dashboard submits source manifests such as `go.mod` or `package.json`.
2. The parser extracts dependency names and versions.
3. The ingestion layer fetches library documentation, extracts relevant content, converts it to markdown, and embeds it.
4. The vector store indexes chunks with dependency metadata.
5. The context builder retrieves version-aware chunks.
6. The prompt builder creates a review prompt.
7. The LLM client sends the prompt to a local model through Ollama or vLLM.
8. The dashboard displays review results for manual validation.

The current implementation is a scaffold. External integrations are represented by in-memory or stub implementations so the project can compile before infrastructure is wired in.

