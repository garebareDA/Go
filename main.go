package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
)

var localAddr = "0.0.0.0:25565"
var remoteAddr = "192.168.1.201:25565"

func proxyConn(conn *net.TCPConn) {
	rAddr, err := net.ResolveTCPAddr("tcp", remoteAddr)
	if err != nil {
		panic(err)
	}

	rConn, err := net.DialTCP("tcp", nil, rAddr)
	if err != nil {
		panic(err)
	}
	defer rConn.Close()

	for {
		data := make([]byte, 256)
		n, err := conn.Read(data)
		if err != nil {
			panic(err)
		}
		if _, err := rConn.Write(data[:n]); err != nil {
			panic(err)
		} else {
			log.Printf("sent:\n%v", hex.Dump(data[:n]))
		}

		rData := make([]byte, 256)
		rn, rerr := rConn.Read(rData)
		if rerr != nil {
			panic(err)
		}
		if _, err := conn.Write(rData[:rn]); err != nil {
			panic(err)
		} else {
			log.Printf("received:\n%v", hex.Dump(rData[:rn]))
		}
	}

	// if _, err := rConn.Write(buf.Bytes()); err != nil {
	// 	panic(err)
	// }
	// log.Printf("sent:\n%v", hex.Dump(buf.Bytes()))

	// for {
	// 	data := make([]byte, 1024)
	// 	n, err := rConn.Read(data)
	// 	if err != nil {
	// 		if err != io.EOF {
	// 			panic(err)
	// 		} else {
	// 			log.Printf("received err: %v", err)
	// 		}
	// 	}
	// 	log.Printf("received:\n%v", hex.Dump(data[:n]))
	// }
}

func handleConn(in <-chan *net.TCPConn, out chan<- *net.TCPConn) {
	for conn := range in {
		proxyConn(conn)
		out <- conn
	}
}

func closeConn(in <-chan *net.TCPConn) {
	for conn := range in {
		conn.Close()
	}
}

func main() {
	fmt.Printf("Listening: %v\nProxying: %v\n\n", localAddr, remoteAddr)

	addr, err := net.ResolveTCPAddr("tcp", localAddr)
	if err != nil {
		panic(err)
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}

	pending, complete := make(chan *net.TCPConn), make(chan *net.TCPConn)

	for i := 0; i < 5; i++ {
		go handleConn(pending, complete)
	}
	go closeConn(complete)

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			panic(err)
		}
		pending <- conn
	}
}
