package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func main() {
	fmt.Println("Running boostrap version v0.5")

	if len(os.Args) < 3 {
		fmt.Println("Usage: boostrap [keepalive] [healthprobe] [<bootstrap-command>] ")
		return
	}

	nextArg := 1
	keepAlive := false
	healthProbe := false

	port := ":8080"
	server := &http.Server{
		Addr: port,
	}

	if len(os.Args) >= (nextArg+1) && os.Args[nextArg] == "keepalive" {
		keepAlive = true
		nextArg++
	}

	if len(os.Args) >= (nextArg+1) && os.Args[nextArg] == "healthprobe" {
		healthProbe = true
		nextArg++

		fmt.Printf("Starting healthprobe on port %s...\n", port)

		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "OK")
		})

		go func() {
			if err := http.ListenAndServe(":"+port, nil); err != nil {
				fmt.Printf("Error starting healthprobe server: %s\n", err)
			}
		}()

	}

	if len(os.Args) >= (nextArg + 1) {
		executable := os.Args[nextArg]
		args := os.Args[(nextArg + 1):]
		cmd := exec.Command(os.Args[nextArg], os.Args[(nextArg+1):]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		fmt.Printf("Bootstrapping: %s %v\n", executable, args)
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Execution failed: %v\n", err)
		}
	}

	if keepAlive {
		fmt.Println("keepalive...")
		if !healthProbe {
			for {
				time.Sleep(time.Hour)
			}
		}
	} else if healthProbe {
		fmt.Println("Shutting down the server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 5-second timeout for shutdown
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			fmt.Printf("Healthprobe shutdown failed: %s\n", err)
		} else {
			fmt.Println("Healthprobe shut down gracefully.")
		}
	}

}
