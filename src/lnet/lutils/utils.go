package lutils

import (
	"strconv"
	"sync"
	"time"
	"unsafe"
)

type UUID struct {
	Timestamp int32
	ServerId  int16
	Index     uint16
}

var uuidIndex uint16 = 0
var uuidMutex sync.Mutex

func getUUID() int64 {
	uuidMutex.Lock()
	uuidIndex++
	uuidMutex.Unlock()
	uuid := UUID{int32(time.Now().Unix()), 0, uuidIndex}
	return *((*int64)(unsafe.Pointer(&uuid)))
}

func GetUUIDStr() string {
	return strconv.Itoa(int(getUUID()))
}
