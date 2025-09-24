package healthprobe

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

// Probes struct manages the HTTP server and probe states.
type HealthProbe struct {
	server    *http.Server
	isReady   atomic.Value // Tracks readiness state
	isStartup atomic.Value // Tracks startup state
	isAlive   atomic.Value // Tracks liveness state
}

// New creates a new Probes instance with default states.
func New(port string) *HealthProbe {
	probes := &HealthProbe{
		server: &http.Server{
			Addr: ":" + port,
		},
	}

	// Initialize states
	probes.isReady.Store(false)   // Not ready by default
	probes.isStartup.Store(false) // Not started by default
	probes.isAlive.Store(true)    // Alive by default

	return probes
}

// Start launches the probe server and sets up handlers.
func (p *HealthProbe) Start() {
	// Handlers for probes
	http.HandleFunc("/liveness", p.LivenessHandler)
	http.HandleFunc("/readiness", p.ReadinessHandler)
	http.HandleFunc("/startup", p.StartupHandler)

	// Start the HTTP server
	go func() {
		fmt.Printf("Starting probe server on port %s\n", p.server.Addr)
		if err := p.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error starting probe server: %s\n", err)
		}
	}()
}

// Stop gracefully shuts down the probe server.
func (p *HealthProbe) Stop(timeout time.Duration) error {
	fmt.Println("Shutting down the probe server...")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Mark as not ready and not alive during shutdown
	p.isReady.Store(false)
	p.isAlive.Store(false)

	if err := p.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("probe server shutdown failed: %w", err)
	}

	fmt.Println("Probe server shut down gracefully.")
	return nil
}

// Handlers for each probe type

// livenessHandler handles the /liveness endpoint.
func (p *HealthProbe) LivenessHandler(w http.ResponseWriter, r *http.Request) {
	if p.isAlive.Load().(bool) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Alive")
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Not Alive")
	}
}

// readinessHandler handles the /readiness endpoint.
func (p *HealthProbe) ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	if p.isReady.Load().(bool) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Ready")
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintln(w, "Not Ready")
	}
}

// startupHandler handles the /startup endpoint.
func (p *HealthProbe) StartupHandler(w http.ResponseWriter, r *http.Request) {
	if p.isStartup.Load().(bool) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Started")
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintln(w, "Not Started")
	}
}

// SetReady allows external components to update the readiness state.
func (p *HealthProbe) SetReady(ready bool) {
	fmt.Printf("Probeserver set 'Ready'-status: %v\n", ready)
	p.isReady.Store(ready)
}

// SetAlive allows external components to update the liveness state.
func (p *HealthProbe) SetAlive(alive bool) {
	fmt.Printf("Probeserver set 'Liveness'-status: %v\n", alive)
	p.isAlive.Store(alive)
}

// SetStartup allows external components to update the startup state.
func (p *HealthProbe) SetStartup(started bool) {
	fmt.Printf("Probeserver set 'Startup'-status: %v\n", started)
	p.isStartup.Store(started)
}
