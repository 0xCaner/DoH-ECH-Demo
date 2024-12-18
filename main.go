package main

import (
	"crypto/tls"
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/miekg/dns"
	"io/ioutil"
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
	domain := "0xcaner.top" // 这个可以改成随便一个cloudflare托管的域名都行

	param1 := flag.String("host", "www.discord.com", "请输入你要访问的HOST (示例: www.discord.com)")
	param2 := flag.String("path", "/", "请输入你要访问的PATH (示例: /)")
	showHelp := flag.Bool("h", false, "显示帮助信息")
	flag.Parse()

	realDomain := *param1
	urlPath := *param2

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
					fmt.Println("本次连接的IP: " + ip)
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

	clientConfig := &tls.Config{
		ServerName: realDomain,
	}
	clientConfig.MinVersion = tls.VersionTLS13

	echBytes, err := base64.StdEncoding.DecodeString(echValue)
	if err != nil {
		log.Fatalf("解码Ech失败: %v", err)
	}
	clientConfig.EncryptedClientHelloConfigList = echBytes

	transport := &http.Transport{
		DialTLS: func(network, addr string) (net.Conn, error) {
			conn, err := net.DialTimeout(network, ip+":443", 10*time.Second)
			if err != nil {
				return nil, err
			}

			tlsConn := tls.Client(conn, clientConfig)

			err = tlsConn.Handshake()
			if err != nil {
				return nil, err
			}

			return tlsConn, nil
		},
		TLSClientConfig: clientConfig,
	}

	client := &http.Client{
		Transport: transport,
	}

	url := fmt.Sprintf("https://%s"+urlPath, ip)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("创建请求失败: %v", err)
	}

	req.Host = realDomain

	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("读取响应失败: %v", err)
	}

	fmt.Printf("响应内容: %s\n", body)
}
