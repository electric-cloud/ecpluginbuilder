all: mac windows linux

mac:
	GOOS=darwin GOARCH=386 go build -o bin/darwin_i686/ecpluginbuilder
	GOOS=darwin GOARCH=amd64 go build -o bin/darwin_x86_64/ecpluginbuilder

windows:
	GOOS=windows GOARCH=386 go build -o bin/windows_i686/ecpluginbuilder.exe
	GOOS=windows GOARCH=amd64 go build -o bin/windows_x86_64/ecpluginbuilder.exe

linux:
	GOOS=linux GOARCH=386 go build -o bin/linux_i686/ecpluginbuilder
	GOOS=linux GOARCH=amd64 go build -o bin/linux_x86_64/ecpluginbuilder

