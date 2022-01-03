package main

import "go-projects/spi/api"

type TestImpl struct {
	api.Test
}

// Name 这里不能按照规范实现接口，否则找不到方法
func Name(param string) string {
	return "Hello SPI param:" + param
}
