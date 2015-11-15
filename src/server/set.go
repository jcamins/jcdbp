package main

func commandSet(key string, val string, connectionChannel chan bool) (success bool) {
    writeChan <- WritePacket{key, val, connectionChannel}
    select {
    case msg := <-connectionChannel:
        return msg
    }
}
