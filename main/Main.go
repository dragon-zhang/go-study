package main

import (
	"fmt"
	"go-projects/demos"
)

// GOPROXY=https://proxy.golang.com.cn
func main() {
	demos.HelloWorld()
	demos.IfElse()
	demos.For()
	demos.Switch(1)
	demos.TestNewClass()
	fmt.Println(demos.Function(1, 2, 3))
	fmt.Println("before defer " + demos.Name)
	fmt.Println(demos.Defer())
	fmt.Println("after defer " + demos.Name)
	demos.Reflect(1, "2", int64(3))
	demos.Goroutine()
	demos.GoInstall()
	//demos.Web()
	demos.Json()
	demos.Channel()
	demos.WaitGroup()
	demos.ThreadId()
	demos.MmapRead("before")
	demos.MmapWrite()
	demos.MmapRead("after")
	demos.Tag()
	demos.RealCoroutine()
}
