package main

// Name 这里不能实现接口，否则找不到方法
func Name(param string) string {
	return "HelloSPI param:" + param
}
