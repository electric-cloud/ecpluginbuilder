package sources

import (
    "fmt"
    "os"
    "path"
    "path/filepath"
)

func CreateBuildDirectory(pluginDirectory string) {
    buildDirectory := path.Join(pluginDirectory, "build")

    err := os.Mkdir(buildDirectory, os.ModePerm)
    if err != nil {
        // panic(err)
        fmt.Println(err)
    }

    dirSuffix := "dsl"
    currentDir := path.Join( pluginDirectory, dirSuffix )
    // TODO create path
    filepath.Walk(currentDir, func(currentFilePath string, fi os.FileInfo, _ error) (err error) {
        relativePath, _ := filepath.Rel(currentDir, currentFilePath)
        newFilePath := path.Join(buildDirectory, dirSuffix, relativePath)

        if fi.IsDir() {
            fmt.Println(relativePath)
            if _, err := os.Stat(newFilePath); os.IsNotExist(err) {
                err = os.Mkdir(newFilePath, os.ModePerm)
            }

            // err = os.Mkdir(path, os.ModePerm)
        }
        return err
        // err := CopyFile( path, newFilePath)
        // if err != nil {
        //     fmt.Println(err)
        // }
        // return nil
    })
}
