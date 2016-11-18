package main

import (
    "fmt"
    // "os"
    // "io"
    // "path/filepath"
    // "path"
    "pluginwiz/sources"
    "pluginwiz/ioutil"
    "strings"
    "pluginwiz/params"
)


func updateVersion( oldVersion string ) string {
    parts := strings.Split(oldVersion, ".")
    if len(parts) == 3 {
        parts = append(parts, "99")
    } else if len(parts) == 4 {
        parts[3] = "100"
    }
    return strings.Join(parts, ".")
}

func main() {
    p := ioutil.GetParameters()
    fmt.Printf("%+v\n", p)
    folders := []string{"t"}
    _ = folders

    // Read all the possible sources of input
    args := params.GetCommandLineArguments()
    pluginDir, err := params.GetPluginDirectory(args)
    pluginXml, err := params.ReadPluginXML(pluginDir)
    if err != nil {
        panic(err)
    }

    version, err := params.GetPluginVersion(args, pluginXml)
    if err != nil {
        panic(err)
    }
    fmt.Println("Plugin version is " + version)

    name, err := params.GetPluginName(args, pluginXml)
    if err != nil {
        panic(err)
    }

    fmt.Println("Plugin name is " + name)
    placeholders := make(map[string]string)
    placeholders["%PLUGIN_KEY%"] = name
    placeholders["%PLUGIN_VERSION%"] = version
    placeholders["%PLUGIN_NAME%"] = name + "-" + version

    sources.CreateBuildTree(pluginDir, folders, placeholders)
}
