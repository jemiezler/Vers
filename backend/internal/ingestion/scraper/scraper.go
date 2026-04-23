package scraper

import "github.com/jemiezler/Vers/backend/internal/ingestion/fetcher"

func ExtractRelevant(doc fetcher.Document) string {
	return doc.Content
}
