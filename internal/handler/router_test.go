package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	appErrors "url-audit-service/internal/errors"
	"url-audit-service/internal/model"

	"github.com/stretchr/testify/assert"
)

type mockAuditService struct {
	auditFunc func(ctx context.Context, url string) (*model.AuditResult, error)
}

func (m *mockAuditService) AuditURL(ctx context.Context, url string) (*model.AuditResult, error) {
	return m.auditFunc(ctx, url)
}

func (m *mockAuditService) GetAuditHistory(ctx context.Context, url string) (*model.AuditResult, error) {
	return nil, nil
}

type mockLimiter struct {
	allowFunc func(ctx context.Context, key string) (bool, error)
}

func (ml *mockLimiter) Allow(ctx context.Context, key string) (bool, error) {
	if ml.allowFunc != nil {
		return ml.allowFunc(ctx, key)
	}
	return true, nil
}

func TestHealthCheck(t *testing.T) {
	// Setup dependencies
	auditSvc := &mockAuditService{}
	auditHandler := NewAuditHandler(auditSvc)
	healthHandler := NewHealthHandler()
	lim := &mockLimiter{}

	router := NewRouter(auditHandler, healthHandler, lim)

	// Perform request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var responseMap map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseMap)
	assert.NoError(t, err)

	assert.True(t, responseMap["success"].(bool))
	dataMap := responseMap["data"].(map[string]interface{})
	assert.Equal(t, "UP", dataMap["status"])
}

func TestAuditEndpoint_Success(t *testing.T) {
	// Setup service mock
	expectedResult := &model.AuditResult{
		URL:            "https://google.com",
		Reachable:      true,
		StatusCode:     http.StatusOK,
		ResponseTimeMs: 120,
		ContentType:    "text/html",
		ContentLength:  1500,
		Title:          "Google",
		Cached:         false,
		CheckedAt:      time.Now(),
	}

	auditSvc := &mockAuditService{
		auditFunc: func(ctx context.Context, url string) (*model.AuditResult, error) {
			return expectedResult, nil
		},
	}
	auditHandler := NewAuditHandler(auditSvc)
	healthHandler := NewHealthHandler()
	lim := &mockLimiter{}

	router := NewRouter(auditHandler, healthHandler, lim)

	// Perform request with a valid external URL
	payload, _ := json.Marshal(model.AuditRequest{URL: "https://google.com"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/audit", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var result model.AuditResult
	err := json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err)

	assert.Equal(t, expectedResult.URL, result.URL)
	assert.Equal(t, expectedResult.Title, result.Title)
	assert.Equal(t, expectedResult.StatusCode, result.StatusCode)
}

func TestAuditEndpoint_Rejected(t *testing.T) {
	auditSvc := &mockAuditService{}
	auditHandler := NewAuditHandler(auditSvc)
	healthHandler := NewHealthHandler()
	lim := &mockLimiter{}

	router := NewRouter(auditHandler, healthHandler, lim)

	// Perform request with private/loopback URL
	payload, _ := json.Marshal(model.AuditRequest{URL: "http://127.0.0.1"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/audit", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var responseMap map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseMap)
	assert.NoError(t, err)

	assert.NotEmpty(t, responseMap["requestId"])
	errObj := responseMap["error"].(map[string]interface{})
	assert.Equal(t, "INVALID_URL", errObj["code"])
	assert.Contains(t, errObj["message"].(string), "loopback")
}

func TestAuditEndpoint_RateLimited(t *testing.T) {
	auditSvc := &mockAuditService{}
	auditHandler := NewAuditHandler(auditSvc)
	healthHandler := NewHealthHandler()

	// Rate limiter rejects requests
	lim := &mockLimiter{
		allowFunc: func(ctx context.Context, key string) (bool, error) {
			return false, nil
		},
	}

	router := NewRouter(auditHandler, healthHandler, lim)

	payload, _ := json.Marshal(model.AuditRequest{URL: "https://google.com"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/audit", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusTooManyRequests, w.Code)

	var responseMap map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseMap)
	assert.NoError(t, err)

	assert.NotEmpty(t, responseMap["requestId"])
	errObj := responseMap["error"].(map[string]interface{})
	assert.Equal(t, "RATE_LIMIT_EXCEEDED", errObj["code"])
	assert.Contains(t, errObj["message"].(string), "rate limit exceeded")
}

func TestAuditEndpoint_Errors(t *testing.T) {
	healthHandler := NewHealthHandler()
	lim := &mockLimiter{}

	t.Run("Bad JSON payload", func(t *testing.T) {
		auditSvc := &mockAuditService{}
		auditHandler := NewAuditHandler(auditSvc)
		router := NewRouter(auditHandler, healthHandler, lim)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/audit", bytes.NewBuffer([]byte("{invalid-json}")))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	tests := []struct {
		name          string
		srvError      error
		expectCode    int
		expectCodeStr string
	}{
		{
			name:          "Invalid URL service error",
			srvError:      appErrors.ErrInvalidURL,
			expectCode:    http.StatusBadRequest,
			expectCodeStr: "INVALID_URL",
		},
		{
			name:          "Timeout service error",
			srvError:      appErrors.ErrURLTimeout,
			expectCode:    http.StatusGatewayTimeout,
			expectCodeStr: "TIMEOUT",
		},
		{
			name:          "Connection failed service error",
			srvError:      appErrors.ErrConnectionFailed,
			expectCode:    http.StatusBadGateway,
			expectCodeStr: "CONNECTION_FAILED",
		},
		{
			name:          "Unexpected internal service error",
			srvError:      errors.New("something went wrong"),
			expectCode:    http.StatusInternalServerError,
			expectCodeStr: "INTERNAL_ERROR",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			auditSvc := &mockAuditService{
				auditFunc: func(ctx context.Context, url string) (*model.AuditResult, error) {
					return nil, tc.srvError
				},
			}
			auditHandler := NewAuditHandler(auditSvc)
			router := NewRouter(auditHandler, healthHandler, lim)

			payload, _ := json.Marshal(model.AuditRequest{URL: "https://google.com"})
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/v1/audit", bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectCode, w.Code)
			var responseMap map[string]interface{}
			_ = json.Unmarshal(w.Body.Bytes(), &responseMap)
			errObj := responseMap["error"].(map[string]interface{})
			assert.Equal(t, tc.expectCodeStr, errObj["code"])
		})
	}
}

func TestAuditEndpoint_History(t *testing.T) {
	auditSvc := &mockAuditService{}
	auditHandler := NewAuditHandler(auditSvc)
	healthHandler := NewHealthHandler()
	lim := &mockLimiter{}

	router := NewRouter(auditHandler, healthHandler, lim)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/audit/history", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "History endpoint placeholder")
}
