package stdhttpclient

import (
	"net/http"
	"net/http/httptest"
	"sync"
)

var (
	// Global test server to avoid port exhaustion
	globalTestServer     *httptest.Server
	globalTestServerOnce sync.Once
)

// getGlobalTestServer returns a shared test server to avoid port exhaustion
func getGlobalTestServer() *httptest.Server {
	globalTestServerOnce.Do(func() {
		// Create a more robust server that can handle high load
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// Simple, fast response for all requests
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})

		globalTestServer = httptest.NewServer(mux)
	})
	return globalTestServer
}
