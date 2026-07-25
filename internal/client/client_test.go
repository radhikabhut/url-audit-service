package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAuditClientSuccess(t *testing.T) {
	// 1. Start mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<html><head><title>Success Title</title></head><body>Content</body></html>`))
	}))
	defer ts.Close()

	// 2. Init client
	cli := NewAuditClient(2 * time.Second)

	// 3. Audit
	res, err := cli.AuditURL(context.Background(), ts.URL)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	assert.Equal(t, ts.URL, res.URL)
	assert.True(t, res.Reachable)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "text/html", res.ContentType)
	assert.Equal(t, "Success Title", res.Title)
	assert.True(t, res.ResponseTimeMs >= 0)
}
