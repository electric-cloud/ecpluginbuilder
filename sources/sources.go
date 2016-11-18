package sources

import (
    "fmt"
    "os"
    "path"
    "path/filepath"
    "pluginwiz/utils"
    "strings"
)

func CreateBuildTree(
    pluginDirectory string,
    subfolders []string,
    placeholders map[string]string) (err error) {

    buildDirectory, err := createBuildDirectory(pluginDirectory)
    if err != nil {
        panic(err)
    }

    for _, folder := range subfolders {
        currentDir := path.Join(pluginDirectory, folder)
        filepath.Walk(currentDir, func(p string, fi os.FileInfo, _ error ) (err error) {
            relativePath, _ := filepath.Rel(currentDir, p)
            newFilePath := path.Join(buildDirectory, folder, relativePath)
            if fi.IsDir() {
                if _, err := os.Stat(newFilePath); os.IsNotExist(err) {
                    err = os.Mkdir(newFilePath, os.ModePerm)
                }
            } else {
                if needToProcessPlaceholders(folder) {
                    replacer := func(content string) string {
                        result := content
                        // Replace all placeholders
                        for placeholder, value := range placeholders {
                            result = strings.Replace(result, placeholder, value, -1)
                        }
                        return result
                    }
                    utils.CopyFileWithContentReplace(p, newFilePath, replacer)
                } else {
                    utils.CopyFile(p, newFilePath)
                }
            }
            return
        })
    }

    fmt.Print("")

    return
}

func needToProcessPlaceholders(folder string) bool {
    // TODO return with more elaborate solution
    if folder == "lib" {
        return false
    } else {
        return true
    }
}

func createBuildDirectory(pluginDirectory string) (buildDirectory string, err error) {
    buildDirectory = path.Join(pluginDirectory, "build")
    if _, err = os.Stat(buildDirectory); os.IsNotExist(err) {
        err = os.Mkdir(buildDirectory, os.ModePerm)
    }
    return
}
