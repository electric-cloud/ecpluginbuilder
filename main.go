package main

import (
    "fmt"
    "pluginwiz/sources"
    "pluginwiz/params"
    "pluginwiz/packer"
)


func Build() {
    // Read all the possible sources of input
    args := params.GetCommandLineArguments()
    pluginDir, err := params.GetPluginDirectory(args)
    pluginXml, err := params.ReadPluginXML(pluginDir)

    if err != nil {
        fmt.Println(err)
        return
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

    projectName := name + "-" + version
    placeholders["%PLUGIN_NAME%"] = projectName

    folders := params.GetFoldersToPack(args)
    fmt.Println(folders)

    buildDirectory, err := sources.CreateBuildTree(pluginDir, folders, projectName, placeholders)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(buildDirectory)
    err = sources.UpdatePluginXML(pluginDir, buildDirectory, version)

    archiveFilename, err := packer.PackZip(folders, buildDirectory, name, version)
    fmt.Println("Build archive: " + archiveFilename)
    fmt.Println("Success!")
}

func main() {
    Build()
}
