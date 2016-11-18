package utils

import (
    "os"
    "io"
    "fmt"
    "bytes"
)


func CopyFileWithContentReplace(src, dst string, replacer FileContentReplacer) (err error) {
    err = copyFileInit(src, dst)
    err = copyFileContents(src, dst, replacer)
    return
}

func CopyFile(src, dst string) (err error) {
    err = copyFileInit(src, dst)
    err = copyFileContents(src, dst, nil)
    return
}

func copyFileInit( src, dst string) (err error) {
    sfi, err := os.Stat(src)
    if err != nil {
        return
    }
    if !sfi.Mode().IsRegular() {
        // cannot copy non-regular files (e.g., directories,
        // symlinks, devices, etc.)
        return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
    }
    dfi, err := os.Stat(dst)
    if err != nil {
        if !os.IsNotExist(err) {
            return
        }
    } else {
        if !(dfi.Mode().IsRegular()) {
            return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
        }
        if os.SameFile(sfi, dfi) {
            return
        }
    }
    return
}


type FileContentReplacer func(content string) (result string)

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string, replacer FileContentReplacer) (err error) {
    in, err := os.Open(src)
    if err != nil {
        return
    }
    defer in.Close()
    out, err := os.Create(dst)
    if err != nil {
        return
    }

    buf := bytes.NewBuffer(nil)
    io.Copy(buf, in)

    if replacer != nil {
        content := string(buf.Bytes())
        content = replacer(content)
        buf = bytes.NewBufferString(content)
    }

    defer func() {
        cerr := out.Close()
        if err == nil {
            err = cerr
        }
    }()
    if _, err = io.Copy(out, buf); err != nil {
        return
    }
    err = out.Sync()
    return
}

