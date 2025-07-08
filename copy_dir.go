package main

import (
    "fmt"
    "io"
    "os"
    "os/exec"
    "path/filepath"
    "time"
)

func copyDir(src, dst string) error {
    return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        relPath, _ := filepath.Rel(src, path)
        destPath := filepath.Join(dst, relPath)

        if info.IsDir() {
            return os.MkdirAll(destPath, info.Mode())
        }

        srcFile, err := os.Open(path)
        if err != nil {
            return err
        }
        defer srcFile.Close()

        destFile, err := os.Create(destPath)
        if err != nil {
            return err
        }
        defer destFile.Close()

        _, err = io.Copy(destFile, srcFile)
        return err
    })
}

func main() {
    fmt.Println("Running copy_dir version v0.4")

    if len(os.Args) < 3 {
        fmt.Println("Usage: copy_dir <source_directory> <target_directory> [keepalive] [<bootstrap-command>] ")
        return
    }

    sourceDir := os.Args[1]
    destinationDir := os.Args[2]

    nextArg := 3

    keepAlive := false
    
    if len(os.Args) >= (nextArg+1) && os.Args[nextArg] == "keepalive" {
        keepAlive = true
        nextArg++
    }
    
    if len(os.Args) >= (nextArg+1) {
        executable := os.Args[nextArg]   
	args := os.Args[(nextArg+1):]   
    	cmd := exec.Command(os.Args[nextArg], os.Args[(nextArg+1):]...)
    	cmd.Stdout = os.Stdout
    	cmd.Stderr = os.Stderr

        fmt.Printf("Bootstrapping: %s %v\n", executable, args)
        if err := cmd.Run(); err != nil {
		    fmt.Fprintf(os.Stderr, "Execution failed: %v\n", err)
	    }
    }
    _, err := os.Stat(destinationDir)
    if os.IsNotExist(err) {    
        err := copyDir(sourceDir, destinationDir)
        if err != nil {
            fmt.Println("Error:", err)
        } else {
            fmt.Println("Directory copied successfully!")
        }
    } else {
        fmt.Println("Sourcedirectory already exists. Skip copying!")
    }
    
    if keepAlive {
        fmt.Println("keepalive...")
        for { 
            time.Sleep(time.Hour)
        }
    }
}
