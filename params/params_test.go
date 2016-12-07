package params

import (
    "fmt"
    "testing"
    "os"
    "reflect"
)


func TestGetCommandLineArguments(t *testing.T) {
    programName := os.Args[0]
    args := []string{programName, "-plugin-version", "1.0.0", "-plugin-name", "EC-Test", "-folder", "f1", "-folder", "f2", "-folder", "f3, f4"}
    os.Args = args
    arguments := GetCommandLineArguments()
    fmt.Println(arguments)
    if arguments.Name != "EC-Test" {
        t.Error("Plugin name argument failed")
    }
    folders := []string{"f1", "f2", "f3", "f4"}
    if !reflect.DeepEqual(folders, arguments.Folders) {
        t.Error("Plugin folder argument failed")
    }
}

func TestGetPluginDirectory(t *testing.T) {
    var args CommandLineArguments
    args.Directory = "dir"

    dir, _ := GetPluginDirectory(args)
    if dir != "dir" {
        t.Error("Cannot get plugin directory")
    }
}

func TestGetPluginVersion(t *testing.T) {
    var args CommandLineArguments
    args.BuildNumber = "2"
    var plugin PluginXML
    plugin.Version = "1.0.0"

    version, _ := GetPluginVersion(args, plugin)
    if version != "1.0.0.2" {
        t.Error("Cannot get plugin version")
    }
}

func TestGetPluginName(t *testing.T) {
    var args CommandLineArguments
    var plugin PluginXML
    args.Name = "EC-Test"
    name, _ := GetPluginName(args, plugin)
    if name != "EC-Test" {
        t.Error("Cannot get plugin name")
    }
}

func TestGetFoldersToPack(t *testing.T) {
    var args CommandLineArguments
    args.Folders = []string{"f1"}

    folders := GetFoldersToPack(args)
    if !reflect.DeepEqual(folders, args.Folders) {
        t.Error("CAnnot get folders to pack")
    }
}
