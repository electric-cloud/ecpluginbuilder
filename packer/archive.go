package packer

import (
    "os"
    "time"
    "archive/zip"
    "path/filepath"
    "path"
    "strings"
    "io"
    "github.com/electric-cloud/ecpluginbuilder/utils"
)


func PackZipVersioned( folders []string, buildDir, name, version string) (string, error) {
    return PackArchiveVersioned(folders, buildDir, name, version, "zip")
}

func PackZipUnversioned(folders []string, buildDir, name string) (string, error) {
    return PackArchiveUnversioned(folders, buildDir, name, "zip")
}

func PackJarVersioned(folders []string, buildDir, name, version string) (string, error) {
    return PackArchiveVersioned(folders, buildDir, name, version, "jar")
}

func PackJarUnversioned(folders []string, buildDir, name string) (string, error) {
    return PackArchiveUnversioned(folders, buildDir, name, "jar")
}

func pack(folders []string, buildDir, target string) (err error) {
    var fullPathFolders []string
    for _, folder := range folders {
        var fullPath string
        if filepath.IsAbs(folder) {
            fullPath = folder
        } else {
            fullPath = path.Join(buildDir, folder)
        }
        var exists bool
        exists, err = utils.FolderExists(fullPath)
        if err != nil {
            return
        }
        if exists {
            fullPathFolders = append(fullPathFolders, fullPath)
        }
    }
    err = Zip(target, fullPathFolders)
    if err != nil {
        return
    }
    return
}



func PackArchiveUnversioned(folders []string, buildDir, name, extension string) (string, error) {
    target := path.Join(buildDir, name + "." + extension)
    err := pack(folders, buildDir, target)
    return target, err
}

func PackArchiveVersioned( folders []string, buildDir, name, version, extension string) (string, error) {
    var target string
    archiveFilename := name + "-" + version + "." + extension
    target = path.Join(buildDir, archiveFilename)
    err := pack(folders, buildDir, target)
    return target, err
}

func PackDependencies(folderName, pluginDir, buildDir string) (archiveName string, err error) {
    now := time.Now()
    archiveName = path.Join(buildDir, folderName + "-" + now.Format("20060102150405") + ".zip")
    folders := make([]string, 1)
    folders[0] = path.Join(pluginDir, folderName)
    err = pack(folders, buildDir, archiveName)
    return
}

// Taken from http://blog.ralch.com/tutorial/golang-working-with-zip/
func Zip(target string, sources []string) error {
    // fmt.Println(target)
    // fmt.Println(sources)
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
                // Windows may have a mixture of shashes (/ and \)
                var trimmed =  strings.TrimPrefix(filepath.ToSlash(p), filepath.ToSlash(source))
				// Because commander needs /
                header.Name = strings.Replace(filepath.Join(baseDir, trimmed), "\\", "/", -1)
				//fmt.Println(header.Name)
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
