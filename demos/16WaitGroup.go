package demos

import (
	"fmt"
	"sync"
	"time"
)

func WaitGroup() {
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		fmt.Println("1 goroutine sleep ...")
		time.Sleep(1 * time.Second)
		fmt.Println("1 goroutine exit ...")
	}()
	go func() {
		defer wg.Done()
		fmt.Println("2 goroutine sleep ...")
		time.Sleep(2 * time.Second)
		fmt.Println("2 goroutine exit ...")
	}()

	fmt.Println("waiting for all goroutine ")
	wg.Wait()
	fmt.Println("All goroutines finished!")
}
