package main

import (
    "fmt"
    "net"
    "os"
    "strings"
    "strconv"
)

var data = make(map[string]string)

func main() {
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

    length, err := conn.Read(buf)
    if err != nil {
        conn.Close()
        fmt.Println("Error reading: ", err.Error())
        return
    }

    str := string(buf[:length])
    tokens := strings.Split(str, "\r\n")
    command := tokens[2]
    args := make([]string, (len(tokens) - 3) / 2)
    for ptr := 0; ptr < len(args); ptr++ {
        args[ptr] = tokens[4 + ptr * 2]
    }
    switch command {
    case "SET":
        data[args[0]] = args[1]
        conn.Write([]byte("+OK\r\n"))
        break;
    case "GET":
        conn.Write([]byte("$"))
        conn.Write([]byte(strconv.Itoa(len(data[args[0]]))))
        conn.Write([]byte("\r\n"))
        conn.Write([]byte(data[args[0]]))
        conn.Write([]byte("\r\n"))
        break;
    }
    conn.Close()
}

/*func main() {
    http.Handle("/key/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        key := r.URL.Path[len("/key/"):]
        switch r.Method {
            case "GET":
                fmt.Fprintf(w, data[key])
                break
            case "POST":
                body, _ := ioutil.ReadAll(r.Body)
                data[key] = string(body)
        }
    }))

    http.ListenAndServe(":9000", nil);
}*/

