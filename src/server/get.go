package main

func commandGet(key string) (val string) {
    dataLock.RLock()
    val = data[key]
    dataLock.RUnlock()
    return val
}
