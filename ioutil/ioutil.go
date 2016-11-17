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


// func CopyFile(src, dst string) (err error) {
//     sfi, err := os.Stat(src)
//     if err != nil {
//         return
//     }
//     if !sfi.Mode().IsRegular() {
//         // cannot copy non-regular files (e.g., directories,
//         // symlinks, devices, etc.)
//         return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
//     }
//     dfi, err := os.Stat(dst)
//     if err != nil {
//         if !os.IsNotExist(err) {
//             return
//         }
//     } else {
//         if !(dfi.Mode().IsRegular()) {
//             return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
//         }
//         if os.SameFile(sfi, dfi) {
//             return
//         }
//     }
//     err = copyFileContents(src, dst)
//     return
// }


// // copyFileContents copies the contents of the file named src to the file named
// // by dst. The file will be created if it does not already exist. If the
// // destination file exists, all it's contents will be replaced by the contents
// // of the source file.
// func copyFileContents(src, dst string) (err error) {
//     in, err := os.Open(src)
//     if err != nil {
//         return
//     }
//     defer in.Close()
//     out, err := os.Create(dst)
//     if err != nil {
//         return
//     }
//     defer func() {
//         cerr := out.Close()
//         if err == nil {
//             err = cerr
//         }
//     }()
//     if _, err = io.Copy(out, in); err != nil {
//         return
//     }
//     err = out.Sync()
//     return
// }
