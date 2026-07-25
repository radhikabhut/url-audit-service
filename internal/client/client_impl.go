package client

import (
	"context"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
	appErrors "url-audit-service/internal/errors"
	"url-audit-service/internal/model"
)

var titleRegex = regexp.MustCompile(`(?is)<title>(.*?)</title>`)

type httpClientImpl struct {
	client *http.Client
}

func NewAuditClient(timeout time.Duration) AuditClient {
	return &httpClientImpl{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (h *httpClientImpl) AuditURL(ctx context.Context, targetURL string) (*model.AuditResult, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
	if err != nil {
		return nil, appErrors.ErrInvalidURL
	}

	// User-Agent to prevent getting blocked by basic scrapers block
	req.Header.Set("User-Agent", "URL-Audit-Service/1.0")

	startTime := time.Now()
	resp, err := h.client.Do(req)
	duration := time.Since(startTime).Milliseconds()

	if err != nil {
		var netErr interface{ Timeout() bool }
		if errors.As(err, &netErr) && netErr.Timeout() {
			return nil, appErrors.ErrURLTimeout
		}
		return nil, appErrors.ErrConnectionFailed
	}
	defer resp.Body.Close()

	// Limit reader to prevent memory exhaustion (e.g., 2MB max)
	limitReader := io.LimitReader(resp.Body, 2*1024*1024)
	bodyBytes, err := io.ReadAll(limitReader)
	if err != nil {
		return nil, appErrors.ErrInternal
	}

	contentLength := resp.ContentLength
	if contentLength <= 0 {
		contentLength = int64(len(bodyBytes))
	}

	contentType := resp.Header.Get("Content-Type")
	// Clean content type if it contains charset
	if idx := strings.Index(contentType, ";"); idx != -1 {
		contentType = strings.TrimSpace(contentType[:idx])
	}

	title := extractTitle(bodyBytes)

	return &model.AuditResult{
		URL:            targetURL,
		Reachable:      true,
		StatusCode:     resp.StatusCode,
		ResponseTimeMs: duration,
		ContentType:    contentType,
		ContentLength:  contentLength,
		Title:          title,
		Cached:         false,
		CheckedAt:      time.Now(),
	}, nil
}

func extractTitle(body []byte) string {
	matches := titleRegex.FindSubmatch(body)
	if len(matches) > 1 {
		// Clean title (remove whitespace and unescape simple XML/HTML elements if needed)
		title := strings.TrimSpace(string(matches[1]))
		// Basic HTML unescape for common ones
		title = strings.ReplaceAll(title, "&amp;", "&")
		title = strings.ReplaceAll(title, "&lt;", "<")
		title = strings.ReplaceAll(title, "&gt;", ">")
		title = strings.ReplaceAll(title, "&quot;", "\"")
		title = strings.ReplaceAll(title, "&#39;", "'")
		return title
	}
	return ""
}
