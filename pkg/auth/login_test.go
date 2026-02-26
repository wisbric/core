package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

type fakeRateLimiter struct {
	result *RateLimitResult
}

func (f *fakeRateLimiter) Check(_ context.Context, _ string) (*RateLimitResult, error) {
	return f.result, nil
}

func (f *fakeRateLimiter) Record(_ context.Context, _ string) error { return nil }
func (f *fakeRateLimiter) Reset(_ context.Context, _ string) error  { return nil }

func TestLoginHandler_RateLimited(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	rl := &fakeRateLimiter{
		result: &RateLimitResult{
			Allowed: false,
			RetryAt: time.Now().Add(2 * time.Second),
		},
	}
	handler := NewLoginHandler(nil, nil, logger, false, rl)

	body, err := json.Marshal(LoginRequest{Email: "user@example.com", Password: "password"})
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	w := httptest.NewRecorder()
	handler.HandleLogin(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusTooManyRequests)
	}

	var resp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decoding response: %v", err)
	}
	if resp["error"] != "rate_limited" {
		t.Fatalf("error = %v, want %q", resp["error"], "rate_limited")
	}
}
