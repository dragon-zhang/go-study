package impl1

import (
	"fmt"
	"go-projects/dubbo-go-spi/example"
)

type TestImpl1 struct {
	example.Test
}

// CallFunc 这里不能按照规范实现接口，否则找不到方法
func CallFunc() {
	fmt.Println("Hello1 SPI")
}
