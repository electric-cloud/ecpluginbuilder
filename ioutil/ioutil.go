package ioutil

import (
    "os"
    "path/filepath"
)

func GetParameters() map[string]string {
    m := make(map[string]string)
    m["pluginVersion"] = "1.0.0.1"
    m["pluginName"] = "Test"

    currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
        panic(err)
    }
    _ = currentDir
    m["pluginDir"] = "/Users/imago/Documents/ecloud/latest_commander/EC-AmazonECS/"

    directoriesToCopy := []string{"dsl", "lib"}
    // m["directoriesToCopy"] = directoriesToCopy
    _ = directoriesToCopy

    return m
}

