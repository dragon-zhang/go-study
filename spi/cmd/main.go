package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"plugin"
)

func main() {
	//todo 后续可以抽象出类似java spi的配置文件
	//	   要求：1.能根据环境加载本地dylib;
	//			2.如果是go源文件，编译出3环境通用的dylib
	soPath := "spi.api.Test/TestImpl.so"
	sourcePath := "spi/impl/TestImpl.go"
	BuildDynamicLibrary(soPath, sourcePath)
	plug, err := plugin.Open(soPath)
	if err != nil {
		panic(err)
	}
	var methodName = "Name"
	function, err := plug.Lookup(methodName)
	if err != nil {
		panic("cannot find func:" + methodName)
	}
	result := function.(func(string) string)("test")
	fmt.Println(result)
}

// BuildDynamicLibrary 手动构建动态so库
func BuildDynamicLibrary(soPath string, sourcePath string) {
	//交叉编译
	exec.Command(
		"set",
		"GOOS", "linux",
		"GOARCH", "x86_64,amd64",
	).Run()
	log.Println("build dynamic library successfully !")
	cmd := exec.Command(
		"go", "build",
		"-buildmode", "plugin",
		"-o", soPath,
		sourcePath,
	)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		log.Fatalf("Error building plugin: %s\nOutput:\n%s", err, out.String())
	}
}
