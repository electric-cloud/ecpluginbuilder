package main

import (
    "fmt"
    // "os"
    // "io"
    // "path/filepath"
    // "path"
    "pluginwiz/sources"
    "pluginwiz/ioutil"
)



func main() {
    params := ioutil.GetParameters()
    fmt.Printf("%+v\n", params)
    sources.CreateBuildDirectory(params["pluginDir"])
}
