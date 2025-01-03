package main

import (
	"crypto/tls"
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/miekg/dns"
	"golang.org/x/net/context"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func createDNSQuery(domain string, qtype uint16) ([]byte, error) {
	msg := dns.Msg{}
	msg.SetQuestion(dns.Fqdn(domain), qtype)
	return msg.Pack()
}

func sendDNSRequest(query []byte) ([]byte, error) {
	url := fmt.Sprintf("https://dns.alidns.com/dns-query?dns=%s", base64.RawURLEncoding.EncodeToString(query))

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content := make([]byte, resp.ContentLength)
	_, err = resp.Body.Read(content)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func main() {
	param1 := flag.String("domain", "0xcaner.top", "指定ECH的来源域名")
	param2 := flag.String("cdnip", "", "指定要访问的CDN IP")
	param3 := flag.String("host", "wss.0xcaner.top", "请输入你要访问的HOST (示例: www.discord.com)")
	param4 := flag.String("path", "/", "请输入你要访问的PATH (示例: /)")
	showHelp := flag.Bool("h", false, "显示帮助信息")
	flag.Parse()

	domain := *param1

	realDomain := *param3
	urlPath := *param4

	if *showHelp {
		flag.Usage()
		return
	}

	queryType := dns.TypeHTTPS

	query, err := createDNSQuery(domain, queryType)
	if err != nil {
		log.Fatal("Error creating DNS query:", err)
	}

	content, err := sendDNSRequest(query)
	if err != nil {
		log.Fatal("Error sending DNS request:", err)
	}

	msg := new(dns.Msg)
	err = msg.Unpack(content)
	if err != nil {
		log.Fatalf("Error unpacking DNS response: %v\n", err)
	}

	fmt.Printf("Query Name: %s\n", msg.Question[0].Name)

	echValue := ""
	ip := ""

	for _, ans := range msg.Answer {
		fmt.Println("Answer:", ans.String())
		if httpsRecord, ok := ans.(*dns.HTTPS); ok {

			for _, i := range httpsRecord.Value {
				switch v := i.(type) {
				case *dns.SVCBECHConfig:
					echValue = base64.StdEncoding.EncodeToString(v.ECH)
					fmt.Println("本次连接的ECH: " + v.String())
					break
				case *dns.SVCBIPv4Hint:
					ip = v.Hint[0].String()
				}

			}
		} else {
			fmt.Println("Not an HTTPS record")
		}
	}

	if ip == "" || echValue == "" {
		fmt.Println("DNS 数据错误，请重试或者域名不支持ECH功能")
		os.Exit(-1)
	}

	if *param2 != "" {
		ip = *param2
	}

	fmt.Println("本次连接的IP: " + ip)
	echBytes, err := base64.StdEncoding.DecodeString(echValue)

	if err != nil {
		log.Fatalf("解码Ech失败: %v", err)
	}

	tlsConfig := &tls.Config{
		MinVersion:                     tls.VersionTLS13,
		ServerName:                     realDomain,
		EncryptedClientHelloConfigList: echBytes,
	}

	addr := fmt.Sprintf("wss://%s:%d%s", realDomain, 443, urlPath)

	dialer := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
		NetDialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			conn, err := net.DialTimeout(network, fmt.Sprintf("%s:%d", ip, 443), 10*time.Second)
			if err != nil {
				return nil, err
			}

			tlsConn := tls.Client(conn, tlsConfig)

			err = tlsConn.Handshake()
			if err != nil {
				return nil, err
			}

			return tlsConn, nil
		},
	}
	headers := http.Header{}
	headers.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
	wssConn, _, err := dialer.Dial(addr, headers)

	if err != nil {
		log.Printf("请求失败: %v\n", err)
		log.Fatalf("可能该CDN IP无法访问，请使用-cdnip参数指定其它的CDN IP")
	}
	defer wssConn.Close()

	sendData := "Hello World!"
	log.Println("发送数据: ", sendData)

	err = wssConn.WriteMessage(websocket.TextMessage, []byte(sendData))

	if err != nil {
		log.Fatalf("发送数据失败: %v", err)
	}

	_, data, err := wssConn.ReadMessage()

	if err != nil {
		log.Fatalf("读取响应失败: %v", err)
	}

	fmt.Printf("响应内容: %s\n", data)
}
