package main

import (
	"encoding/gob"
	"log"
	"os"
	"sync"
	"time"
)

var data = make(map[string]string)
var changeCount = 0
var dataLock sync.RWMutex
var serializeChan = make(chan bool)

func WriteThread(writeChan <-chan WritePacket) {
	for msg := range writeChan {
		dataLock.Lock()
		changeCount++
		data[msg.key] = msg.val
		dataLock.Unlock()
		if changeCount >= 100 && len(serializeChan) == 0 {
			select {
			case serializeChan <- true:
			default:
			}
		}
		msg.notify <- true
	}
}

func DiskThread() {
	for {
		select {
		case <-serializeChan:
			WriteToDisk()
		case <-time.After(time.Second * 10):
			WriteToDisk()
		}
	}
}

func WriteToDisk() {
	dataLock.Lock()
	if changeCount == 0 {
		dataLock.Unlock()
		return
	}
	copy := make(map[string]string, len(data))
	for k, v := range data {
		copy[k] = v
	}
	changeCount = 0
	dataLock.Unlock()
	file, err := os.Create("datafile.gob")
	if err == nil {
		enc := gob.NewEncoder(file)
		err = enc.Encode(copy)
		if err != nil {
			log.Print("Unable to encode:", err)
		}
		file.Close()
	}
}

func ReadFromDisk() {
	file, err := os.Open("datafile.gob")
	if err == nil {
		dec := gob.NewDecoder(file)
		err = dec.Decode(&data)
		file.Close()
		if err != nil {
			log.Fatal("Unable to decode:", err)
		}
	}
}
