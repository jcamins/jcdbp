package main

import (
    "fmt"
    "net"
    "os"
    "strings"
    "strconv"
)

func startListener() {
    listener, err := net.Listen("tcp", "localhost:6379")
    if err != nil {
        fmt.Println("Unable to open port: ", err.Error())
        os.Exit(1)
    }
    defer listener.Close()
    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Error with connection: ", err.Error())
        } else {
            go handleRequest(conn)
        }
    }
}

func handleRequest(conn net.Conn) {
    buf := make([]byte, 1024)
    channel := make(chan bool)

    for {
        length, err := conn.Read(buf)
        if err != nil {
            conn.Close()
            if (err.Error() != "EOF") {
                fmt.Println("Error reading: ", err.Error())
            }
            break
        }

        str := string(buf[:length])
        tokens := strings.Split(str, "\r\n")
        command := strings.ToUpper(tokens[2])
        args := make([]string, (len(tokens) - 3) / 2)
        for ptr := 0; ptr < len(args); ptr++ {
            args[ptr] = tokens[4 + ptr * 2]
        }
        switch command {
        case "SET":
            if commandSet(args[0], args[1], channel) {
                conn.Write([]byte("+OK\r\n"))
            } else {
                conn.Write([]byte("-ERR\r\n"))
            }
        case "GET":
            val := commandGet(args[0])
            conn.Write([]byte("$"))
            conn.Write([]byte(strconv.Itoa(len(val))))
            conn.Write([]byte("\r\n"))
            conn.Write([]byte(val))
            conn.Write([]byte("\r\n"))
        case "QUIT":
            conn.Write([]byte("+OK\r\n"))
            conn.Close()
            break
        }
    }
}
