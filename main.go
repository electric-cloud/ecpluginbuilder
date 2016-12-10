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
    placeholders["@PLUGIN_KEY@"] = name
    placeholders["@PLUGIN_VERSION@"] = version

    projectName := name + "-" + version
    placeholders["@PLUGIN_NAME@"] = projectName

    folders, err := params.GetFoldersToPack(args)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Folders to pack: %v", folders)

    buildDirectory, err := sources.CreateBuildTree(pluginDir, folders, projectName, placeholders, args)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println("Build directory is " + buildDirectory)
    err = sources.UpdatePluginXML(pluginDir, buildDirectory, version)

    archiveFilenameVersioned, err := packer.PackZipVersioned(folders, buildDirectory, name, version)
    fmt.Println("Build archive: " + archiveFilenameVersioned)
    archiveFilenameUnversioned, err := packer.PackZipUnversioned(folders, buildDirectory, name)
    fmt.Println("Build archive: " + archiveFilenameUnversioned)
    fmt.Println("Success!")
}

func main() {
    Build()
}
