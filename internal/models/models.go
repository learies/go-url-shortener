package models

type Storage struct {
	Id          string `db:"id" json:"id"`
	ShortURL    string `db:"short_url" json:"short_url"`
	OriginalURL string `db:"original_url" json:"original_url"`
	UserID      string `db:"user_id" json:"user_id"`
	DeletedFlag bool   `db:"is_deleted" json:"is_deleted"`
}

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}

type BatchURLRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchURLResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type BatchURLWrite struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
	OriginalURL   string `json:"original_url"`
	UserID        string `json:"user_id"`
}

type URL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type UserURL struct {
	UserID   string `json:"user_id"`
	ShortURL string `json:"short_url"`
}

type ShortURLs struct {
	ShortURLs []string `json:"short_urls"`
}
