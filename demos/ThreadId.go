package demos

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

// GetThreadId 获取线程id
func GetThreadId() int {
	return os.Getpid()
}

func ThreadId() {
	//让协程只由1个线程调度
	runtime.GOMAXPROCS(1)
	go func() {
		fmt.Printf("coroutine %v\n", GetThreadId())
	}()
	fmt.Printf("main %v\n", GetThreadId())
	time.Sleep(time.Second)
}
