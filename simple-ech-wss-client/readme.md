# 用法说明

ECH Websockets服务端和客户端的简单示例。

服务端使用Python启动普通的wss服务端，接着利用Cloudflare接入CDN，通过Cloudflare的ECH服务进行测试。



## 服务端

1. python启动服务端，会在443端口启动一个websocket服务端

```shell
root@VM-0-6-ubuntu:~# python3 sslapp.py 
WebSocket server started on wss://0.0.0.0:443
```

2. 接入CloudFlare CDN
3. 测试开放情况，因为这是websocket服务端，https直接访问报错是正常的

```shell
直接在机器上使用IP:端口访问:
root@VM-0-1-ubuntu:~# curl https://127.0.0.1 -k
Failed to open a WebSocket connection: empty Connection header.

You cannot access a WebSocket server directly with a browser. You need a WebSocket client.

通过CloudFlare CDN的域名访问:
root@VM-0-6-ubuntu:~# curl https://test.0xcaner.top 
curl: (92) HTTP/2 stream 0 was not closed cleanly: PROTOCOL_ERROR (err 1)
```





## 客户端

编译

```shell
./build.sh
```

各参数含义如下:

```shell
➜  simple-ech-wss-client git:(main) ✗ ./DoH-ECH-wss-darwin-amd64 -h
Usage of ./DoH-ECH-wss-darwin-amd64:
  -cdnip string
        指定要访问的CDN IP
  -domain string
        指定ECH的来源域名 (default "0xcaner.top")
  -h    显示帮助信息
  -host string
        请输入你要访问的HOST (示例: www.discord.com) (default "wss.0xcaner.top")
  -path string
        请输入你要访问的PATH (示例: /) (default "/")
```

使用-host指定你的CDN域名，当默认获取到的CDN IP无法连接时，需要指定cdnip参数:

```shell
➜  simple-wss-client git:(main) ✗ go run ./ -host=test.0xcaner.top -cdnip=104.16.92.14
Query Name: 0xcaner.top.
Answer: 0xcaner.top.    1       IN      HTTPS   1 . alpn="h3,h2" ipv4hint="104.21.16.1,104.21.32.1,104.21.48.1,104.21.64.1,104.21.80.1,104.21.96.1,104.21.112.1" ech="AEX+DQBBmQAgACBlefMTE62+KY35J9If5YS9yxi0SY4nWB+CdL1FdH9/EQAEAAEAAQASY2xvdWRmbGFyZS1lY2guY29tAAA=" ipv6hint="2606:4700:3030::6815:1001,2606:4700:3030::6815:2001,2606:4700:3030::6815:3001,2606:4700:3030::6815:4001,2606:4700:3030::6815:5001,2606:4700:3030::6815:6001,2606:4700:3030::6815:7001"
本次连接的ECH: AEX+DQBBmQAgACBlefMTE62+KY35J9If5YS9yxi0SY4nWB+CdL1FdH9/EQAEAAEAAQASY2xvdWRmbGFyZS1lY2guY29tAAA=
本次连接的IP: 104.16.92.14
2025/01/03 15:05:28 发送数据:  Hello World!
响应内容: Echo: Hello World!
```

