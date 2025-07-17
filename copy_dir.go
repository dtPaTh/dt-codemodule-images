package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "github.com/dtPaTh/dt-codemodule-images/keepalive"
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
        fmt.Println("Usage: copy_dir <source_directory> <target_directory> [keepalive]")
        return
    }

    sourceDir := os.Args[1]
    destinationDir := os.Args[2]

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
    
    if len(os.Args) >= 4 && os.Args[3] == "keepalive" {
        stopChan := make(chan struct{})
		keepalive.KeepAlive(stopChan)
    }
}
