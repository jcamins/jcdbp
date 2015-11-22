package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var data = make(map[string]string)
var changeCount = 0
var dataLock sync.RWMutex
var serializeChan = make(chan bool)
var loggerChan = make(chan WritePacket, 100)

func WriteThread(writeChan <-chan WritePacket) {
	for msg := range writeChan {
		dataLock.Lock()
		changeCount++
		loggerChan <- msg
		data[msg.Key] = msg.Val
		dataLock.Unlock()
		if changeCount >= 1000000 && len(serializeChan) == 0 {
			select {
			case serializeChan <- true:
			default:
			}
		}
		msg.notify <- true
	}
}

func OpenLog(counter int8) (*gob.Encoder, *os.File) {
	logFile, _ := os.Create("data/datafile.log." + fmt.Sprintf("%03d", counter))
	return gob.NewEncoder(logFile), logFile
}

func LoggerThread() {
	var fileCounter int8 = 0
	var recordCounter = 0
	for fileCounter = 0; ; fileCounter++ {
		_, err := os.Stat("data/datafile.log." + fmt.Sprintf("%03d", fileCounter))
		if err != nil {
			break
		}
	}
	logger, file := OpenLog(fileCounter)
	for {
		select {
		case packet := <-loggerChan:
			if packet.Key == "" {
				file.Close()
				for ; fileCounter >= 0; fileCounter-- {
					os.Remove("data/datafile.log." + fmt.Sprintf("%03d", fileCounter))
				}
				fileCounter = 0
				logger, file = OpenLog(fileCounter)
				recordCounter = 0
			} else {
				//log.Print("logging value")
				err := logger.Encode(packet)
				if err != nil {
					log.Print("Error:", err)
				}
				recordCounter++
			}
		}
	}
}

func DiskThread() {
	for {
		select {
		case <-serializeChan:
			WriteToDisk()
		case <-time.After(time.Second * 60):
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
	loggerChan <- WritePacket{"", "", nil}
	dataLock.Unlock()
	file, err := os.Create("data/datafile.gob")
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
	file, err := os.Open("data/datafile.gob")
	if err == nil {
		dec := gob.NewDecoder(file)
		err = dec.Decode(&data)
		file.Close()
		if err != nil {
			log.Fatal("Unable to decode:", err)
		}
	}
	var packet WritePacket
	for ii := 0; ii < 128; ii++ {
		file, err := os.Open("data/datafile.log." + fmt.Sprintf("%03d", ii))
		if err != nil {
			break
		} else {
			dec := gob.NewDecoder(file)
			for dec.Decode(&packet) == nil {
				data[packet.Key] = packet.Val
			}
			file.Close()
		}
	}
}
