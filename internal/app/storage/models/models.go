package models

import "encoding/json"

type URL struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

type URLBatch struct {
	CorrelationID string `json:"correlation_id"`
	URL
}

func (u *URLBatch) MarshalJSON() ([]byte, error) {
	// API response does not need OriginalURL
	data := map[string]string{
		"correlation_id": u.CorrelationID,
		"short_url":      u.ShortURL,
	}
	return json.Marshal(data)
}
