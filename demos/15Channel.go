package demos

import (
	"fmt"
	"strconv"
)

func Channel() {
	test := make(chan string)
	go func() {
		//send data
		fmt.Println("before write to channel")
		for i := 0; i < 10; i++ {
			test <- strconv.Itoa(i)
		}
		fmt.Println("after write to channel")
		//写完后释放channel，不然其他地方会一直等待
		close(test)
	}()
	for data := range test {
		fmt.Println("read data:" + data + " from channel")
	}
	BufferedChannel()
}

func BufferedChannel() {
	//因为带缓冲区，才能这样写，不然会写阻塞
	ch := make(chan string, 3)
	ch <- "1"
	ch <- "2"
	ch <- "3"
	fmt.Println(<-ch)
	fmt.Println(<-ch)
	fmt.Println(<-ch)
}
