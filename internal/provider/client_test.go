package provider

import (
	"context"
	"net/http"
	"net"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestServer(t *testing.T, handler http.Handler) *httptest.Server {
	t.Helper()
	ln, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to create listener: %v", err)
	}
	server := httptest.NewUnstartedServer(handler)
	server.Listener = ln
	server.Start()
	return server
}

func newTestClient(t *testing.T, server *httptest.Server) *apiClient {
	t.Helper()
	client := newClient("secret", server.URL)
	client.httpClient = server.Client()
	return client
}

func TestCreateUser_Success(t *testing.T) {
	server := newTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/v1/domain/example/user") {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"sub":"abc-123"}`))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, server)

	sub, err := client.CreateUser(context.Background(), "example", "alice", "Alice Example")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sub != "abc-123" {
		t.Fatalf("expected sub abc-123, got %s", sub)
	}
}

func TestCreateUser_ConflictAdoptsExisting(t *testing.T) {
	server := newTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/v1/domain/example/user"):
			w.WriteHeader(http.StatusConflict)
			_, _ = w.Write([]byte(`{"error":"conflict"}`))
		case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/v1/domain/example/user"):
			if got := r.URL.Query().Get("search"); got != "alice" {
				t.Fatalf("expected search query alice, got %s", got)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`[{"sub":"existing-1","preferredUsername":"alice","identityProvider":"foo"}]`))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, server)

	sub, err := client.CreateUser(context.Background(), "example", "alice", "Alice Example")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sub != "existing-1" {
		t.Fatalf("expected existing user id, got %s", sub)
	}
}

func TestFindUser_NotFoundReturnsNil(t *testing.T) {
	server := newTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, server)

	user, err := client.FindUser(context.Background(), "example", "missing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user != nil {
		t.Fatalf("expected nil user when not found")
	}
}

func TestDeleteUser_SucceedsOnNoContent(t *testing.T) {
	server := newTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, server)

	if err := client.DeleteUser(context.Background(), "example", "abc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteUser_NotFoundIsTolerated(t *testing.T) {
	server := newTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, server)

	if err := client.DeleteUser(context.Background(), "example", "abc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDoRequest_PropagatesAPIError(t *testing.T) {
	server := newTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`boom`))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, server)

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("unexpected request error: %v", err)
	}
	_, status, reqErr := client.doRequest(req)
	if status != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", status)
	}
	if reqErr == nil {
		t.Fatalf("expected error, got nil")
	}
	if _, ok := reqErr.(*apiError); !ok {
		t.Fatalf("expected apiError, got %T", reqErr)
	}
}
