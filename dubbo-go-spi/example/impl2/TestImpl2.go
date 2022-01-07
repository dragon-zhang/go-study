package impl2

import (
	"fmt"
	"go-projects/dubbo-go-spi/example"
)

type TestImpl2 struct {
	example.Test
}

// CallFunc 这里不能按照规范实现接口，否则找不到方法
func CallFunc() {
	fmt.Println("Hello2 SPI")
}
