package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func startListener() {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:6379")
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Println("Unable to open port: ", err.Error())
		os.Exit(1)
	}
	defer listener.Close()
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("Error with connection: ", err.Error())
		} else {
			go handleRequest(conn)
		}
	}
}

func handleRequest(conn *net.TCPConn) {
	buf := make([]byte, 1024)
	channel := make(chan bool)

	for {
		length, err := conn.Read(buf)
		if err != nil {
			conn.Close()
			if err.Error() != "EOF" {
				fmt.Println("Error reading: ", err.Error())
			}
			break
		}

		str := string(buf[:length])
		tokens := strings.Split(str, "\r\n")
		command := strings.ToUpper(tokens[2])
		args := make([]string, (len(tokens)-3)/2)
		for ptr := 0; ptr < len(args); ptr++ {
			args[ptr] = tokens[4+ptr*2]
		}
		switch command {
		case "SET":
			if commandSet(args[0], args[1], channel) {
				conn.Write([]byte("+OK\r\n"))
			} else {
				conn.Write([]byte("-ERR\r\n"))
			}
		case "GET":
			var wbuf bytes.Buffer
			val := commandGet(args[0])
			wbuf.Write([]byte("$"))
			wbuf.Write([]byte(strconv.Itoa(len(val))))
			wbuf.Write([]byte("\r\n"))
			wbuf.Write([]byte(val))
			wbuf.Write([]byte("\r\n"))
			wbuf.WriteTo(conn)
		case "QUIT":
			conn.Write([]byte("+OK\r\n"))
			conn.Close()
			break
		case "DIE":
			conn.Write([]byte("+OK\r\n"))
			os.Exit(0)
		}
	}
}
