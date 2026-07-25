package model

import "time"

type AuditRequest struct {
	URL string `json:"url" binding:"required,url"`
}

type AuditResult struct {
	URL            string    `json:"url"`
	Reachable      bool      `json:"reachable"`
	StatusCode     int       `json:"statusCode"`
	ResponseTimeMs int64     `json:"responseTimeMs"`
	ContentType    string    `json:"contentType"`
	ContentLength  int64     `json:"contentLength"`
	Title          string    `json:"title"`
	Cached         bool      `json:"cached"`
	CheckedAt      time.Time `json:"checkedAt"`
}
