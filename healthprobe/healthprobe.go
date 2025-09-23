package healthprobe

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type HealthProbe struct {
	server *http.Server
}

// New creates a new HealthProbe instance.
func New(port string) *HealthProbe {
	return &HealthProbe{
		server: &http.Server{
			Addr: ":" + port,
		},
	}
}

// Start launches the health probe server on the specified port.
func (hp *HealthProbe) Start() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	go func() {
		fmt.Printf("Starting healthprobe on port %s/health\n", hp.server.Addr)
		if err := hp.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error starting healthprobe server: %s\n", err)
		}
	}()
}

// Stop gracefully shuts down the health probe server.
func (hp *HealthProbe) Stop(timeout time.Duration) error {
	fmt.Println("Shutting down the healthprobe server...")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := hp.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("healthprobe shutdown failed: %w", err)
	}

	fmt.Println("Healthprobe shut down gracefully.")
	return nil
}
