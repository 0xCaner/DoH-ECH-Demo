# 为 Linux 64-bit 编译
GOOS=linux GOARCH=amd64 go build -o DoH2ECH-demo-linux-amd64

# 为 Windows 64-bit 编译
GOOS=windows GOARCH=amd64 go build -o DoH2ECH-demo-windows-amd64.exe

# 为 macOS 64-bit 编译
GOOS=darwin GOARCH=amd64 go build -o DoH2ECH-demo-darwin-amd64