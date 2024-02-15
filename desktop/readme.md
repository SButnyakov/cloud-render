desktop

Windows build: go build main.go windows/render-amd64.exe GOOS=windows GOARCH=amd64 
Mac build: GOOS=darwin GOARCH=amd64 go build -o macOS/render-amd64-darwin main.go
