package main

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"bytes"
	"strings"
	"io"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// 启动代理服务器，监听 9999 端口。
	// 这里没有写ip地址，默认在所有ip地址上进行监听。
	l, err := net.Listen("tcp", "9999")
	if err != nil {
		log.Panic(err)
	}

	// 不断接受请求
	for {
		client, err := l.Accept()
		if err != nil {
			log.Panic(err)
		}

		// 使用一个 goroutine 来处理请求。
		go handleClientRequest(client)
	}
}

func handleClientRequest(client net.Conn) {
	if client == nil {
		return
	}
	defer client.Close()

	var b [1024]byte
	n, err := client.Read(b[:])
	if err != nil {
		log.Println(err)
		return
	}

	var method, host, address string
	_, _ = fmt.Scanf(string(b[:bytes.IndexByte(b[:], '\n')]), "%s%s", &method, &host)
	hostPortURL, err := url.Parse(host)
	if err != nil {
		log.Println(err)
		return
	}

	if hostPortURL.Opaque == "443" {
		address = hostPortURL.Scheme + ":443"
	} else {
		if strings.Index(hostPortURL.Host, ":") == -1 {
			address = hostPortURL.Host + ":80"
		} else {
			address = hostPortURL.Host
		}
	}

	// 获得了请求的host和port，就开始拨号。
	server, err := net.Dial("tcp", address)
	if err != nil {
		log.Println(err)
		return
	}

	if  method == "CONNECT" {
		_, _ = fmt.Fprint(client, "HTTP/1.1 200 Connection established\r\n\r\n")
	} else {
		_, _ = server.Write(b[:n])
	}

	go io.Copy(server, client)
	_, _ = io.Copy(client, server)
}
