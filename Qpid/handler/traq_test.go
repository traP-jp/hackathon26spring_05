package handler

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/traPtitech/go-traq"
)

func TestTraqAPIBaseURL(t *testing.T) {
	tests := []struct {
		name string
		host string
		want string
	}{
		{
			name: "host only",
			host: "https://q.trap.jp",
			want: "https://q.trap.jp/api/v3",
		},
		{
			name: "already api v3",
			host: "https://q.trap.jp/api/v3",
			want: "https://q.trap.jp/api/v3",
		},
		{
			name: "trailing slash",
			host: "https://q.trap.jp/",
			want: "https://q.trap.jp/api/v3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := traqAPIBaseURL(tt.host); got != tt.want {
				t.Fatalf("traqAPIBaseURL(%q) = %q, want %q", tt.host, got, tt.want)
			}
		})
	}
}

func TestGetUserUUIDDoesNotUseNameQuery(t *testing.T) {
	var gotQuery string
	cfg := traq.NewConfiguration()
	cfg.Servers = traq.ServerConfigurations{
		{URL: "https://traq.example/api/v3"},
	}
	cfg.HTTPClient = &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			if r.URL.Path != "/api/v3/users" {
				t.Fatalf("path = %q, want /api/v3/users", r.URL.Path)
			}
			gotQuery = r.URL.RawQuery
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body: io.NopCloser(strings.NewReader(
					`[{"id":"uuid-tidus","name":"tidus","displayName":"Tidus","iconFileId":"icon","bot":false,"state":1,"updatedAt":"2026-06-28T00:00:00Z"}]`,
				)),
			}, nil
		}),
	}

	h := &handler{
		traq: traqClientWithContext{
			client:  traq.NewAPIClient(cfg),
			context: context.Background(),
		},
	}

	uuid, err := h.getUserUUID("tidus")
	if err != nil {
		t.Fatalf("getUserUUID returned error: %v", err)
	}
	if uuid != "uuid-tidus" {
		t.Fatalf("uuid = %q, want %q", uuid, "uuid-tidus")
	}
	if gotQuery != "" {
		t.Fatalf("query = %q, want empty", gotQuery)
	}

	displayName, err := h.fetchUserDisplayName("tidus")
	if err != nil {
		t.Fatalf("fetchUserDisplayName returned error: %v", err)
	}
	if displayName == nil || *displayName != "Tidus" {
		t.Fatalf("displayName = %v, want Tidus", displayName)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}
