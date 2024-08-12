package services

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListmonkClient(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/tx":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"data": true}`))
		case "/api/subscribers":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"data": {"id": 1}}`))
		case "/api/campaigns":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"data": {"id": 1}}`))
		default:
			http.Error(w, "Not found", http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Create a client that uses the mock server URL
	client := &ListmonkClient{
		baseURL: server.URL,
		client:  server.Client(),
	}

	// Test SendTransactionalEmail
	err := client.SendTransactionalEmail(1, "test@example.com", map[string]interface{}{"name": "Test"})
	assert.NoError(t, err)

	// Test ManageSubscriber
	err = client.ManageSubscriber("test@example.com", "Test User", "enabled", []int{1})
	assert.NoError(t, err)

	// Test CreateCampaign
	err = client.CreateCampaign("Test Campaign", "Test Subject", []int{1}, 1, "2023-01-01T00:00:00Z")
	assert.NoError(t, err)
}
