package healthprobe_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dtPaTh/dt-codemodule-images/healthprobe"
)

func TestHealthProbe(t *testing.T) {
	// Create a new HealthProbe instance
	probe := healthprobe.New("8080")

	// Start the probe server in a separate goroutine
	go probe.Start()

	// Ensure the server has time to start
	time.Sleep(100 * time.Millisecond)

	// Test /liveness endpoint
	t.Run("Liveness Probe", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/liveness", nil)
		rec := httptest.NewRecorder()

		// Simulate the liveness handler
		probe.SetAlive(true)
		probe.LivenessHandler(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rec.Code)
		}
		if rec.Body.String() != "Alive\n" {
			t.Errorf("Expected body 'Alive', got %q", rec.Body.String())
		}

		// Mark as not alive and test again
		probe.SetAlive(false)
		rec = httptest.NewRecorder()
		probe.LivenessHandler(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", rec.Code)
		}
		if rec.Body.String() != "Not Alive\n" {
			t.Errorf("Expected body 'Not Alive', got %q", rec.Body.String())
		}
	})

	// Test /readiness endpoint
	t.Run("Readiness Probe", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/readiness", nil)
		rec := httptest.NewRecorder()

		// Simulate the readiness handler
		probe.SetReady(true)
		probe.ReadinessHandler(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rec.Code)
		}
		if rec.Body.String() != "Ready\n" {
			t.Errorf("Expected body 'Ready', got %q", rec.Body.String())
		}

		// Mark as not ready and test again
		probe.SetReady(false)
		rec = httptest.NewRecorder()
		probe.ReadinessHandler(rec, req)

		if rec.Code != http.StatusServiceUnavailable {
			t.Errorf("Expected status 503, got %d", rec.Code)
		}
		if rec.Body.String() != "Not Ready\n" {
			t.Errorf("Expected body 'Not Ready', got %q", rec.Body.String())
		}
	})

	// Test /startup endpoint
	t.Run("Startup Probe", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/startup", nil)
		rec := httptest.NewRecorder()

		// Simulate the startup handler
		probe.SetStartup(true)
		probe.StartupHandler(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rec.Code)
		}
		if rec.Body.String() != "Started\n" {
			t.Errorf("Expected body 'Started', got %q", rec.Body.String())
		}

		// Mark as not started and test again
		probe.SetStartup(false)
		rec = httptest.NewRecorder()
		probe.StartupHandler(rec, req)

		if rec.Code != http.StatusServiceUnavailable {
			t.Errorf("Expected status 503, got %d", rec.Code)
		}
		if rec.Body.String() != "Not Started\n" {
			t.Errorf("Expected body 'Not Started', got %q", rec.Body.String())
		}
	})

	// Stop the server
	probe.Stop(1 * time.Second)
}
