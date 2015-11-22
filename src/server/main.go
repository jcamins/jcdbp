package main

type WritePacket struct {
	Key, Val string
	notify   chan bool
}

var writeChan = make(chan WritePacket)

func main() {
	ReadFromDisk()
	go WriteThread(writeChan)
	go LoggerThread()
	go DiskThread()
	startListener()
}
