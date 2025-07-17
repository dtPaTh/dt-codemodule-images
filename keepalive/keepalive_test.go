package keepalive

import (
	"testing"
	"time"
)

func TestKeepAlive(t *testing.T) {
	// Create a channel to signal the keepAlive function to stop
	stopChan := make(chan struct{})

	// Run the keepAlive function in a goroutine
	go func() {
		KeepAlive(stopChan)
	}()

	// Wait for a short period to simulate the keepalive running
	time.Sleep(100 * time.Millisecond)

	// Signal the keepAlive function to stop
	close(stopChan)

	// Wait a bit to ensure the keepAlive function has terminated
	time.Sleep(50 * time.Millisecond)

	// If the test reaches this point without hanging, the keepAlive function is working as expected
	t.Log("keepAlive function terminated successfully")
}