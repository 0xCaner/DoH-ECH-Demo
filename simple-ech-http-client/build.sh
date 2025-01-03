# 为 Linux 64-bit 编译
GOOS=linux GOARCH=amd64 go build -o DoH-ECH-http-linux-amd64

# 为 Windows 64-bit 编译
GOOS=windows GOARCH=amd64 go build -o DoH-ECH-http-windows-amd64.exe

# 为 macOS 64-bit 编译
GOOS=darwin GOARCH=amd64 go build -o DoH-ECH-http-darwin-amd64