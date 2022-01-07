package impl2

import (
	"fmt"
	"github.com/dragon-zhang/go-study/dubbo-go-spi/example"
)

type TestImpl3 struct {
	example.Test
}

// CallFunc 这里不能按照规范实现接口，否则找不到方法
func CallFunc() {
	fmt.Println("Hello3 SPI")
}
