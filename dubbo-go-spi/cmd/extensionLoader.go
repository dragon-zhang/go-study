package main

import (
	"fmt"
	"go-projects/dubbo-go-spi/api"
)

func main() {
	loader, err := api.GetExtensionLoader("go-projects/dubbo-go-spi/example", "Test", "config")
	fmt.Println(err)
	extension := loader.GetExtension("testImpl1")
	fmt.Println(extension)
	_, err = api.BuildDylib("go-projects/dubbo-go-spi/example/impl3/TestImpl3.go")
	fmt.Println(err)
	extension = loader.GetExtension("testImpl3")
	fmt.Println(extension)
}
