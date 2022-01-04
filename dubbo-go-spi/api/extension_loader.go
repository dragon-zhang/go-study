package api

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"
	"runtime"
	"strings"
)

// interfaces interfaceName : spiName : impl or .so or .dll
var interfaces = make(map[string]map[string]*plugin.Plugin)

// 初始化
func init() {
	err := filepath.Walk("./config",
		func(path string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if f.IsDir() {
				return nil
			}
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			name := file.Name()
			//todo 校验接口的合理性
			fmt.Println(name)
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			extensions := make(map[string]*plugin.Plugin)
			for _, line := range strings.Split(string(data), "\n") {
				split := strings.Split(line, "=")
				spiName := split[0]
				lib := split[1]
				if strings.Contains(lib, ".") {
					lib = BuildDynamicLibrary(lib, lib)
				}
				plug, err := plugin.Open(lib)
				if err != nil {
					return err
				}
				extensions[spiName] = plug
			}
			interfaces[name] = extensions
			return nil
		})
	if err != nil {
		log.Fatalf("Error read extension config: %s\n", err)
	}
}

// GetExtension 获取扩展
func GetExtension(interfaceName string, spiName string) (*plugin.Plugin, error) {
	i := interfaces[interfaceName]
	if i != nil {
		extension := i[spiName]
		if extension != nil {
			return extension, nil
		}
	}
	return nil, errors.New(interfaceName + " extension:" + spiName + " not find")
}

// LookupMethod 根据方法名称查找方法
func LookupMethod(interfaceName string, spiName string, methodName string) (plugin.Symbol, error) {
	extension, _ := GetExtension(interfaceName, spiName)
	if extension != nil {
		return extension.Lookup(methodName)
	}
	return nil, errors.New("method:" + methodName + " can not find in:" + interfaceName + " " + spiName)
}

// BuildDynamicLibrary 构建动态链接库
func BuildDynamicLibrary(targetPath string, sourcePath string) string {
	split := strings.Split(targetPath, ".")
	if runtime.GOOS == "windows" {
		targetPath = split[0] + ".dll"
	} else {
		targetPath = split[0] + ".so"
	}
	cmd := exec.Command(
		"go", "build",
		"-buildmode", "plugin",
		"-o", targetPath,
		sourcePath,
	)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		log.Fatalf("Error building plugin: %s\nOutput:\n%s", err, out.String())
	}
	return targetPath
}
