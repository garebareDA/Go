package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/net/proxy"
)

var (
	proxy_addr  = "0.0.0.0:25565"
	remote_addr = "192.168.1.201:25565"
)

func main() {
	dialer, err := proxy.SOCKS5("tcp", proxy_addr, nil, proxy.Direct)
	if err != nil {
		fmt.Fprintln(os.Stderr, "proxy connection error:", err)
		os.Exit(1)
	}
	conn, err := dialer.Dial("tcp", remote_addr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "remote connection error:", err)
		os.Exit(1)
	}
	defer conn.Close()

	log.Printf("received:\n%v", conn)

	// communicate with remote addr here
}
