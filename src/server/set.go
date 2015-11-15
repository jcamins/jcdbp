package main

func commandSet(key string, val string) {
    writeChan <- WritePacket{key, val}
}
