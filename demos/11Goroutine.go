package demos

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"time"
)

func Goroutine() {
	go fmt.Println("Hello Goroutine !")
	time.Sleep(100 * time.Millisecond)
	fmt.Println("end")
}

// GoroutineId 获取协程id
func GetGoroutineId() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
