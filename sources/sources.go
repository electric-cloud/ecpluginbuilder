package sources

import (
    "os"
    "path"
    "path/filepath"
    "github.com/electric-cloud/ecpluginbuilder/utils"
    "github.com/electric-cloud/ecpluginbuilder/params"
    "strings"
    "io/ioutil"
    "encoding/xml"
    "regexp"
)


func CreateBuildTree(
    pluginDirectory string,
    subfolders []string,
    projectName string,
    placeholders map[string]string, args params.CommandLineArguments) (buildDirectory string, err error) {

    buildDirectory, err = createBuildDirectory(pluginDirectory, args.PreserveBuild)
    if err != nil {
        panic(err)
    }

    for _, folder := range subfolders {
        currentDir := path.Join(pluginDirectory, folder)
        filepath.Walk(currentDir, func(p string, fi os.FileInfo, _ error ) (err error) {
            relativePath, _ := filepath.Rel(currentDir, p)
            newFilePath := path.Join(buildDirectory, folder, relativePath)
            if fi.IsDir() {
                if _, err := os.Stat(newFilePath); os.IsNotExist(err) {
                    err = os.Mkdir(newFilePath, os.ModePerm)
                }
            } else {
                if needToProcessPlaceholders(folder) {
                    replacer := func(content string) string {
                        result := content
                        // Replace all placeholders
                        for placeholder, value := range placeholders {
                            result = strings.Replace(result, placeholder, value, -1)
                        }
                        return result
                    }
                    utils.CopyFileWithContentReplace(p, newFilePath, replacer)
                } else {
                    utils.CopyFile(p, newFilePath)
                }
            }
            return
        })
    }

    err = BuildProjectXML(pluginDirectory, buildDirectory, projectName, placeholders)

    return
}


func UpdatePluginXML(pluginDir, pluginBuild, version string) (err error) {
    pluginXmlPath := path.Join(pluginDir, "META-INF", "plugin.xml")
    pluginXML, err := ioutil.ReadFile(pluginXmlPath)
    if err != nil {
        return
    }
    // Not very flexible, but simple and fast way to do things
    // TODO make is nice
    re := regexp.MustCompile("<version>.+?</version>")
    updated := re.ReplaceAllString(string(pluginXML), "<version>" + version + "</version>")

    destination := path.Join(pluginBuild, "META-INF", "plugin.xml")
    err = ioutil.WriteFile(destination, []byte(updated), os.ModePerm)
    return
}


func BuildProjectXML(pluginDir, pluginBuild, projectName string, placeholders map[string]string) (err error) {
    filename := path.Join(pluginBuild, "META-INF", "project.xml")
    // This one will be processed separately
    ecPerlFilename := path.Join(pluginDir, "ec_setup.pl")
    ecPerl, err := os.Open(ecPerlFilename)
    if err != nil {
        return
    }
    defer ecPerl.Close()
    b, err := ioutil.ReadAll(ecPerl)
    if err != nil {
        return
    }
    escapedCode := string(b)
    // TODO make a separate function

    for placeholder, value := range placeholders {
        escapedCode = strings.Replace(escapedCode, placeholder, value, -1)
    }

    exportedData := &ExportedData{
        XMLName: xml.Name{Local: "exportedData"},
        BuildLabel: "build_3.5_30434_OPT_2010.01.13_07:32:22",
        Version: "39",
        BuildVersion: "3.5.1.30434",
        ExportPath: "/projects/" + projectName,
        Project: Project{
            ProjectName: projectName,
            PropertySheet: PropertySheet{[]Property{
                Property{Value: escapedCode, PropertyName: "ec_setup", Expandable: 0},
            }},
        },
    }

    out, err := xml.MarshalIndent(exportedData, "", "  ")
    if err != nil {
        return
    }

    err = ioutil.WriteFile(filename, out, os.ModePerm)
    if err != nil {
        return
    }
    return nil
}

type ExportedData struct {
    XMLName xml.Name
    BuildLabel string `xml:"buildLabel,attr"`
    BuildVersion string `xml:"buildVersion,attr"`
    Version string `xml:"version,attr"`
    ExportPath string `xml:"exportPath"`
    Project Project `xml:"project"`
}

type Project struct {
    ProjectName string `xml:"projectName"`
    PropertySheet PropertySheet `xml:"propertySheet"`
}

type PropertySheet struct {
    Property []Property `xml:"property"`
}

type Property struct {
    Expandable int `xml:"expandable"`
    PropertyName string `xml:"propertyName"`
    Value string `xml:"value"`
}

func needToProcessPlaceholders(folder string) bool {
    // TODO return with more elaborate solution
    if folder == "lib" {
        return false
    } else {
        return true
    }
}

func createBuildDirectory(pluginDirectory string, preserveBuild bool) (buildDirectory string, err error) {
    buildDirectory = path.Join(pluginDirectory, "build")
    if _, err = os.Stat(buildDirectory); os.IsNotExist(err) {
        err = os.Mkdir(buildDirectory, os.ModePerm)
        return
    } else {
        if preserveBuild {
            return
        }

        err = os.RemoveAll(buildDirectory)
        if err != nil {
            return
        }
        err = os.Mkdir(buildDirectory, os.ModePerm)
        return
        // build directory must be cleaned
    }
    return
}
