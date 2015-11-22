package main

import (
	"sync"
)

var data = make(map[string]string)
var dataLock sync.RWMutex

func WriteThread(writeChan <-chan WritePacket) {
	for msg := range writeChan {
		dataLock.Lock()
		data[msg.key] = msg.val
		dataLock.Unlock()
		msg.notify <- true
	}
}
