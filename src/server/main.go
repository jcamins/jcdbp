package main

type WritePacket struct {
    key, val string
    notify chan bool
}

var writeChan = make(chan WritePacket)

func main() {
    go Writer(writeChan)
    startListener()
}
