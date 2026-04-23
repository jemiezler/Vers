```mermaid
flowchart LR

%% =========================
%% USER INTERFACE
%% =========================
UI[Developer] -->|Upload Code/Manifests| Parser[Code & Dependency Parser]
Dashboard[Manual Review Dashboard]

%% =========================
%% ORCHESTRATION
%% =========================
Parser -->|go.mod / package.json| ContextBuilder[Context Builder]
Parser -->|Library + Version List| ContextBuilder

ContextBuilder -->|Retrieved Docs| PromptBuilder[Prompt Builder]
ContextBuilder -->|Version-filtered Query| VectorDB

%% =========================
%% INGESTION
%% =========================
Docs[Fetch Docs<br/>pkg.go.dev / GitHub] --> Scraper[Targeted Scraper]
Scraper --> Converter[Markdown Converter]
Converter --> Embedder1[Embedding Model (local)]

%% =========================
%% KNOWLEDGE BASE
%% =========================
Embedder1 --> VectorDB[Metadata-aware Vector DB<br/>(Qdrant / Chroma)]
VectorDB -->|Retrieved Chunks| ContextBuilder

%% =========================
%% INFERENCE
%% =========================
PromptBuilder -->|Augmented Prompt| LLM[Gemma 4 (Ollama/vLLM)]
LLM -->|Code Review Result| Dashboard

%% =========================
%% FEEDBACK LOOP
%% =========================
Dashboard --> UI
```

## Project Layout

```text
backend/      Go API and review pipeline modules
frontend/     Manual review dashboard shell
deployments/  Local infrastructure definitions
docs/         Architecture and API notes
testdata/     Sample manifests for parser/review testing
```

## Local Development

Run the backend:

```sh
cd backend
go run ./cmd/api
```

### LLM Configuration

By default the backend uses a stub LLM (`VERS_LLM_PROVIDER=stub`). To use Ollama:

```sh
VERS_LLM_PROVIDER=ollama
VERS_OLLAMA_URL=http://localhost:11434
VERS_OLLAMA_MODEL=gemma3
```

PowerShell example:

```powershell
$env:VERS_LLM_PROVIDER="ollama"
$env:VERS_OLLAMA_URL="http://localhost:11434"
$env:VERS_OLLAMA_MODEL="gemma3"
```

### Docs Configuration

By default the backend uses placeholder docs (`VERS_DOCS_PROVIDER=stub`). To fetch Go docs from pkg.go.dev for `go.mod` dependencies:

```powershell
$env:VERS_DOCS_PROVIDER="pkg_go_dev"
$env:VERS_PKG_GO_DEV_URL="https://pkg.go.dev"
```

Run the frontend:

```sh
cd frontend
npm install
npm run dev
```

The scaffold uses stub document fetching, in-memory vector storage, and a stub LLM client. Replace those modules with Qdrant/Chroma and Ollama/vLLM integrations as the pipeline matures.
