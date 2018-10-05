package sources

import (
    "os"
    "path"
    "path/filepath"
    "github.com/electric-cloud/ecpluginbuilder/utils"
    "github.com/electric-cloud/ecpluginbuilder/params"
    "github.com/electric-cloud/ecpluginbuilder/packer"
    "strings"
    "io/ioutil"
    "encoding/xml"
    "regexp"
    "fmt"
    "encoding/base64"
    "strconv"
    "crypto/md5"
    "encoding/hex"
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
        var exists bool
        exists, err = utils.FolderExists(currentDir)
        _ = exists
        if err != nil {
            return
        }
        if !exists {
            fmt.Println("WARNING: " + currentDir + " does not exist")
            continue
        }
        fmt.Println("Going to folder: " + currentDir)
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

    err = BuildProjectXML(pluginDirectory, buildDirectory, projectName, placeholders, args)

    return
}


func UpdatePluginXML(pluginDir, pluginBuild, version string) (err error) {
    pluginXmlPath := path.Join(pluginBuild, "META-INF", "plugin.xml")
    fmt.Println("Plugin XML path is " + pluginXmlPath)
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


func InjectPromotionIntoSetup(pluginDir, ecSetup string) (injected string, err error) {
    injected = ecSetup
    promoteFileName := path.Join(pluginDir, "dsl", "promote.groovy")
    promoteCode, err := readStringFromFile(promoteFileName)
    if err != nil {
        return
    }
    demoteFileName := path.Join(pluginDir, "dsl", "demote.groovy")
    demoteCode, err := readStringFromFile(demoteFileName)
    if err != nil {
        return
    }

    promotePlaceholder := "# promote.groovy placeholder"
    demotePlaceholder := "# demote.groovy placeholder"

    if !strings.Contains(ecSetup, promotePlaceholder) || !strings.Contains(ecSetup, demotePlaceholder) {
        fmt.Println("WARNING: ec_setup.pl does not contain placeholders for promote.groovy and demote.groovy. See PLUGINWIZ-8 for details.")
        return
    }
    injected = strings.Replace(ecSetup, promotePlaceholder, promoteCode, -1)
    injected = strings.Replace(injected, demotePlaceholder, demoteCode, -1)
    return
}

func packDependencies(pluginDir, pluginBuild string) (base64Depedencies string, err error) {
    libsFolder := path.Join(pluginDir, "lib")
    exists, err := utils.FolderExists(libsFolder)
    if err != nil {
        return
    }

    if !exists {
        fmt.Println("/lib folder does not exist")
        base64Depedencies = ""
        return
    }

    packedFolder, err := packer.PackDependencies("lib", pluginDir, pluginBuild)
	if err != nil {
		return
	}

    fmt.Println("Packed dependencies: " + packedFolder)

    binaryContent, err := ioutil.ReadFile(packedFolder)
    if err != nil {
        return
    }
    err = os.Remove(packedFolder)
    if err != nil {
        return
    }
    base64Depedencies = base64.StdEncoding.EncodeToString(binaryContent)
    return
}


func BuildProjectXML(pluginDir, pluginBuild, projectName string, placeholders map[string]string, args params.CommandLineArguments) (err error) {
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
    escapedCode, err = InjectPromotionIntoSetup(pluginDir, escapedCode)

    if err != nil {
        return
    }

    for placeholder, value := range placeholders {
        escapedCode = strings.Replace(escapedCode, placeholder, value, -1)
    }



    dependencies, err := packDependencies(pluginDir, pluginBuild)
    if err != nil {
        return
    }

    chunkSize := args.DependencyChunkSize
    fmt.Printf("Dependency chunk size is %d\n", chunkSize)
    checksum := getMD5Hash(dependencies)
    chunkedDependencies := chunkString(chunkSize, dependencies)
    var deps string
    for _, chunk := range chunkedDependencies {
        deps += chunk
    }
    chunkedDependenciesStruct := &PropertySheet{[]Property{}}

    for index, chunk := range chunkedDependencies {
        property := Property{
            Value: chunk,
            PropertyName: "ec_dependencyChunk_" + strconv.Itoa(index),
        }
        chunkedDependenciesStruct.Property = append(chunkedDependenciesStruct.Property, property)
    }
    chunkedDependenciesStruct.Property = append(chunkedDependenciesStruct.Property, Property{
        Value: checksum,
        PropertyName: "checksum",
        Description: "MD5 checksum of ZIP archive",
    })

    fmt.Println("Checksum is " + checksum)

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

    if args.PackedLib == false {
        exportedData.Project.PropertySheet.Property = append(exportedData.Project.PropertySheet.Property, Property{PropertySheet: chunkedDependenciesStruct, PropertyName: "ec_groovyDependencies", Description: "Packed .jar dependencies"})
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
    Description string `xml:"description"`
    PropertySheet *PropertySheet `xml:"propertySheet,omitempty"`
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


func readStringFromFile(filename string) (content string, err error) {
    handler, err := os.Open(filename)
    if err != nil {
        return
    }
    defer handler.Close()
    b, err := ioutil.ReadAll(handler)
    if err != nil {
        return
    }
    content = string(b)
    return
}


func chunkString(chunkSize int, s string) (result []string) {
    result = make([]string, 0)
    runes := []rune(s)
    var chunk []rune
    for len(runes) > chunkSize {
        chunk, runes = runes[0:chunkSize], runes[chunkSize:]
        result = append(result, string(chunk))
    }
    result = append(result, string(runes))
    return
}

func getMD5Hash(s string) (result string) {
    hasher := md5.New()
    hasher.Write([]byte(s))
    result = hex.EncodeToString(hasher.Sum(nil))
    return
}
