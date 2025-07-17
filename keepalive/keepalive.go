package keepalive

import (
    "fmt"
    "time"
)


func KeepAlive(stopChan <-chan struct{}) {
	fmt.Println("keepalive...")
	for {
		select {
		case <-stopChan:
			fmt.Println("stop signal received")
			return
		case <-time.After(time.Hour):
		}
	}
}
