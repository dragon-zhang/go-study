package demos

import (
	"fmt"
	"time"
)

func Goroutine() {
	go fmt.Println("Hello Goroutine !")
	time.Sleep(100 * time.Millisecond)
	fmt.Println("end")
}
