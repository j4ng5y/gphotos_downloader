all: test build

.PHONY: test
test:
	@go test -coverprofile=coverage.out ./... -race

.PHONY: build
build:
	@GOOS=linux GOARCH=amd64 go build -o bin/gphotos_downloader_linux_amd64 cmd/gphotos_downloader/gphotos_downloader.go
	@GOOS=windows GOARCH=amd64 go build -o bin/gphotos_downloader_windows_amd64.exe cmd/gphotos_downloader/gphotos_downloader.go
	@GOOS=darwin GOARCH=amd64 go build -o bin/gphotos_downloader_darwin_amd64 cmd/gphotos_downloader/gphotos_downloader.go