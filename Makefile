all:
	go build -o build/darwin/pluginwiz
	GOOS=windows GOARCH=386 go build -o build/windows/pluginwiz
	GOOS=linux GOARCH=386 go build -o build/linux/pluginwiz
