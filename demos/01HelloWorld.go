package demos

import (
	"fmt"
	"runtime"
)

func HelloWorld() {
	fmt.Println(runtime.Version())
	fmt.Println("Hello World !")
}
