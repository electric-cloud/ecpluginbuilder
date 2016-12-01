package packer

import (
    "os"
    "archive/zip"
    // "archive"
    "path/filepath"
    "path"
    "strings"
    "io"
    "fmt"
)


func PackZipVersioned( folders []string, buildDir, name, version string) (string, error) {
    var target string
    fmt.Println()
    archiveFilename := name + "-" + version + ".zip"
    target = path.Join(buildDir, archiveFilename)
    err := pack(folders, buildDir, target)
    return target, err
}

func pack(folders []string, buildDir, target string) (err error) {
    var fullPathFolders []string
    for _, folder := range folders {
        fullPathFolders = append(fullPathFolders, path.Join(buildDir, folder))
    }
    err = Zip(target, fullPathFolders)
    if err != nil {
        return
    }
    return
}

func PackZipUnversioned(folders []string, buildDir, name string) (string, error) {
    target := path.Join(buildDir, name + ".zip")
    err := pack(folders, buildDir, target)
    return target, err
}

// Taken from http://blog.ralch.com/tutorial/golang-working-with-zip/
func Zip(target string, sources []string) error {
    zipfile, err := os.Create(target)
    if err != nil {
        return err
    }
    defer zipfile.Close()

    archive := zip.NewWriter(zipfile)
    defer archive.Close()

    for _, source := range sources {

        info, err := os.Stat(source)
        if err != nil {
            return nil
        }

        var baseDir string
        if info.IsDir() {
            baseDir = filepath.Base(source)
        }

        filepath.Walk(source, func(p string, info os.FileInfo, err error) error {
            if err != nil {
                return err
            }

            header, err := zip.FileInfoHeader(info)
            if err != nil {
                return err
            }

            if baseDir != "" {
                header.Name = filepath.Join(baseDir, strings.TrimPrefix(p, source))
            }

            if info.IsDir() {
                header.Name += "/"
            } else {
                header.Method = zip.Deflate
            }

            writer, err := archive.CreateHeader(header)
            if err != nil {
                return err
            }

            if info.IsDir() {
                return nil
            }

            file, err := os.Open(p)
            if err != nil {
                return err
            }
            defer file.Close()
            _, err = io.Copy(writer, file)
            return err
        })
    }

    return err
}
