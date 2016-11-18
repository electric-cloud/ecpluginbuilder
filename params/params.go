package params

import (
    "flag"
    "fmt"
    "strings"
    "os"
    "path/filepath"
    "io/ioutil"
    "encoding/xml"
)


type PluginXML struct {
    Key string `xml:"key"`
    Version string `xml:"version"`
}

type CommandLineArguments map[string]string


func AugmentPluginVersion(currentVersion, buildNumber string) (string, error) {
    if buildNumber == "" {
        return currentVersion, nil
    }
    parts := strings.Split(currentVersion, ".")
    if len(parts) == 3 {
        parts = append(parts, buildNumber)
    } else if len(parts) == 4 {
        parts[3] = buildNumber
    } else {
        err := fmt.Errorf("Cannot recognize version format %s", currentVersion)
        return currentVersion, err
    }
    newVersion := strings.Join(parts, ".")
    return newVersion, nil
}

func getBuildNumberFromVersion(version string) string {
    parts := strings.Split(version, ".")
    if len(parts) == 4 {
        return parts[3]
    } else {
        return ""
    }
}

func GetPluginDirectory(args CommandLineArguments) (string, error) {
    if args["pluginDir"] != "" {
        return args["pluginDir"], nil
    }
    cwd, err := os.Getwd()
    // TODO maybe we need gitroot here
    return cwd, err
}

func GetPluginVersion(args map[string]string, pluginXML PluginXML) (string, error) {
    var buildNumber string = ""
    if args["build-number"] != "" {
        buildNumber = args["build-number"]
    } else if os.Getenv("BUILD_NUMBER") != "" {
        buildNumber = os.Getenv("BUILD_NUMBER")
    }

    var pluginVersion string
    if args["version"] != "" {
        pluginVersion = args["version"]
    } else {
        pluginVersion = pluginXML.Version
    }

    if pluginVersion == "" {
        return "", fmt.Errorf("Cannot determine plugin version")
    }

    if buildNumber == "" {
        buildNumber = getBuildNumberFromVersion(pluginVersion)
    }

    pluginVersion, err := AugmentPluginVersion(pluginVersion, buildNumber)
    return pluginVersion, err
}

func GetPluginName(args CommandLineArguments, pluginXml PluginXML) (string, error) {
    if args["name"] != "" {
        return args["name"], nil
    }
    if pluginXml.Key != "" {
        return pluginXml.Key, nil
    }
    return "", fmt.Errorf("Cannot determine plugin name")
}

func ReadPluginXML(pluginDirectory string) (p PluginXML, err error) {
    pluginXmlPath := filepath.Join(pluginDirectory, "META-INF", "plugin.xml")
    xmlFile, err := os.Open(pluginXmlPath)
    if err != nil {
        return
    }
    defer xmlFile.Close()

    b, err := ioutil.ReadAll(xmlFile)
    if err != nil {
        return
    }
    xml.Unmarshal(b, &p)
    return
}

func GetFoldersToPack(args CommandLineArguments) []string {
    if args["folders"] != "" {
        folders := strings.Split(args["folders"],",")
        return folders
    } else {
        return []string{"lib", "dsl", "t", "htdocs", "pages", "META-INF"}
    }
}

func GetCommandLineArguments() CommandLineArguments {
    versionPtr := flag.String("plugin-version", "", "Plugin Version")
    namePtr := flag.String("plugin-name", "", "Plugin name")
    pluginDirPtr := flag.String("plugin-dir", "", "Plugin directory")
    foldersToPackPtr := flag.String("folders", "", "Folders to pack")
    flag.Parse()
    m := make(CommandLineArguments)

    m["version"] = *versionPtr
    m["name"] = *namePtr
    m["pluginDir"] = *pluginDirPtr
    m["folders"] = *foldersToPackPtr
    return m
}
