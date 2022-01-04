package impl2

import (
	"go-projects/dubbo-go-spi/example"
)

type TestImpl3 struct {
	example.Test
}

// Name 这里不能按照规范实现接口，否则找不到方法
func Name(param string) string {
	return "Hello2 SPI param:" + param
}
