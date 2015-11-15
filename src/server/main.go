package main

type WritePacket struct {
    key, val string
}

var writeChan = make(chan WritePacket)

func main() {
    go Writer(writeChan)
    startListener()
}
