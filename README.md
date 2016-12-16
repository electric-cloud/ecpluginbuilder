ecpluginbuilder
===============
# A simple build tool for plugins based on PluginWizard

## Install
- Download ecpluginbuilder for your platform from [here](bin/)

## Usage
- ecpluginbuilder -plugin-version 1.0.0.2 -folders t,META-INF,lib,dsl,pages,htdocs -plugin-dir <path-to-plugin>

## Development

### Requirements 
- Go  
- make or emake  

### Setup
- Set the GOPATH environment variable to tell the Go tool where your workspace is located.
- `go get github.com/electric-cloud/ecpluginbuilder`
- `cd $GOPATH/src/github.com/electric-cloud/ecpluginbuilder`
- `emake` or `make`
