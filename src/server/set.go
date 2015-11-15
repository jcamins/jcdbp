package main

func commandSet(key string, val string, channel chan bool) (success bool) {
    writeChan <- WritePacket{key, val, channel}
    return <-channel
}
