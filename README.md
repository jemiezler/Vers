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

Run the frontend:

```sh
cd frontend
npm install
npm run dev
```

The scaffold uses stub document fetching, in-memory vector storage, and a stub LLM client. Replace those modules with Qdrant/Chroma and Ollama/vLLM integrations as the pipeline matures.
