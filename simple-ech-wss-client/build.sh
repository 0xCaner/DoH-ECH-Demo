# 为 Linux 64-bit 编译
GOOS=linux GOARCH=amd64 go build -o DoH-ECH-wss-linux-amd64

# 为 Windows 64-bit 编译
GOOS=windows GOARCH=amd64 go build -o DoH-ECH-wss-windows-amd64.exe

# 为 macOS 64-bit 编译
GOOS=darwin GOARCH=amd64 go build -o DoH-ECH-wss-darwin-amd64