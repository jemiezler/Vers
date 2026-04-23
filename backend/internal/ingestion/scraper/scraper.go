package scraper

import "vers/backend/internal/ingestion/fetcher"

func ExtractRelevant(doc fetcher.Document) string {
	return doc.Content
}
