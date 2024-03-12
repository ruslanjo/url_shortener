package models

import "encoding/json"

type URLBatch struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
	ShortURL      string
}

func (u *URLBatch) MarshalJSON() ([]byte, error) {
	// API response does not need OriginalURL
	data := map[string]string{
		"correlation_id": u.CorrelationID,
		"short_url":      u.ShortURL,
	}
	return json.Marshal(data)
}
