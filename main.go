package main

import (
    "fmt"
    "github.com/electric-cloud/ecpluginbuilder/sources"
    "github.com/electric-cloud/ecpluginbuilder/params"
    "github.com/electric-cloud/ecpluginbuilder/packer"
)

var version string

func Build() {
    // Read all the possible sources of input
    fmt.Println("ecpluginbuilder version is " + version)
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
    fmt.Printf("Folders to pack: %v\n", folders)

    buildDirectory, err := sources.CreateBuildTree(pluginDir, folders, projectName, placeholders, args)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println("Build directory is " + buildDirectory)
    err = sources.UpdatePluginXML(pluginDir, buildDirectory, version)

    var archiveFilenameVersioned, archiveFilenameUnversioned string
    if args.IsJar {
        archiveFilenameVersioned, err = packer.PackJarVersioned(folders, buildDirectory, name, version)
        archiveFilenameUnversioned, err = packer.PackJarUnversioned(folders, buildDirectory, name)
    } else {
        archiveFilenameVersioned, err = packer.PackZipVersioned(folders, buildDirectory, name, version)
        archiveFilenameUnversioned, err = packer.PackZipUnversioned(folders, buildDirectory, name)
    }

    fmt.Println("Build archive: " + archiveFilenameVersioned)
    fmt.Println("Build archive: " + archiveFilenameUnversioned)
    fmt.Println("Success!")
}

func main() {
    Build()
}
