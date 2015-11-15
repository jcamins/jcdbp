package main

var data = make(map[string]string)

func Writer(writeChan <-chan WritePacket) {
    for msg := range writeChan {
        data[msg.key] = msg.val
    }
}
