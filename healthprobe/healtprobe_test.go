package healthprobe

import (
	"io"
	"net/http"
	"testing"
	"time"
)

func TestHealthProbe(t *testing.T) {
	port := "8081"
	hp := New(port)

	// Start the healthprobe server
	hp.Start()
	defer func() {
		if err := hp.Stop(5 * time.Second); err != nil {
			t.Fatalf("Failed to stop healthprobe: %v", err)
		}
	}()

	// Wait briefly to ensure the server has started
	time.Sleep(100 * time.Millisecond)

	// Send a GET request to the /health endpoint
	resp, err := http.Get("http://localhost:" + port + "/health")
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	// Verify the response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	// Verify the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	expectedBody := "OK\n"
	if string(body) != expectedBody {
		t.Errorf("Expected response body %q, got %q", expectedBody, string(body))
	}
}