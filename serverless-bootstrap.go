package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
    "github.com/dtPaTh/dt-codemodule-images/keepalive"
	"github.com/dtPaTh/dt-codemodule-images/healthprobe"

)

func main() {
	fmt.Println("Running serverless-boostrap version v0.6")

	if len(os.Args) < 3 {
		fmt.Println("Usage: serverless-boostrap [--keepalive] [--healthprobe] [<command-to-execute> ...] ")
		return
	}

	nextArg := 1
	keepAliveFlag := false
	healthProbeFlag := false

	port := "8080"
	var hp *healthprobe.HealthProbe

	if len(os.Args) >= (nextArg+1) && os.Args[nextArg] == "--keepalive" {
		keepAliveFlag = true
		nextArg++
	}

	if len(os.Args) >= (nextArg+1) && os.Args[nextArg] == "--healthprobe" {
		healthProbeFlag = true
		nextArg++

		hp = healthprobe.New(port)
		hp.Start()
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

	if keepAliveFlag {
		stopChan := make(chan struct{})
		keepalive.KeepAlive(stopChan)
	} 
	
	if healthProbeFlag  {
		if healthProbeFlag && hp != nil {
			if err := hp.Stop(5 * time.Second); err != nil {
				fmt.Printf("Healthprobe shutdown error: %s\n", err)
			}
		}
	}
		
}
