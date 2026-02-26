package httpserver_test

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wisbric/core/pkg/httpserver"
	"github.com/wisbric/core/pkg/telemetry"
)

func TestHealthz_Always200(t *testing.T) {
	cfg := httpserver.ServerConfig{
		CORSAllowedOrigins: []string{"*"},
	}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	reg := telemetry.NewMetricsRegistry()
	srv := httpserver.NewServer(cfg, logger, nil, nil, reg, nil, nil, nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()

	srv.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}

	var body map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("decoding body: %v", err)
	}
	if body["status"] != "ok" {
		t.Errorf("status = %q, want %q", body["status"], "ok")
	}
}
