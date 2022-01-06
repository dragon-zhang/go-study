package api

import (
	"bytes"
	"errors"
	"github.com/modern-go/reflect2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"
	"reflect"
	"runtime"
	"strings"
)

// SPI配置文件根目录，类似java的 META-INF/services
var configRoot = "./config"

// allExtension interfaceType : spiName : impl or .so or .dll
var allExtension = make(map[reflect2.Type]map[string]*plugin.Plugin)
var cache = make(map[*plugin.Plugin]plugin.Symbol)

// 初始化
func init() {
	err := filepath.Walk(configRoot,
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
			split := strings.Split(file.Name(), "@")
			if len(split) != 2 {
				//非法配置文件
				return nil
			}
			pkg := strings.ReplaceAll(split[0], "#", "/")
			structName := split[1]
			//todo 这里只能映射已发布的struct
			//	   如果有interface必须用struct代替
			theType := reflect2.TypeByPackageName(pkg, structName)
			if theType == nil {
				//struct不存在
				return nil
			}
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			extensions := make(map[string]*plugin.Plugin)
			for _, line := range strings.Split(string(data), "\n") {
				split := strings.Split(line, "=")
				spiName := split[0]
				lib := split[1]
				if strings.LastIndex(lib, ".go") == len(lib)-3 {
					//需要编译
					lib = BuildDynamicLibrary(lib, lib)
				} else {
					//不需要编译
					lib = GetDynamicLibraryName(lib)
				}
				plug, err := plugin.Open(lib)
				if err != nil {
					return err
				}
				extensions[spiName] = plug
			}
			allExtension[theType] = extensions
			return nil
		})
	if err != nil {
		log.Fatalf("Error read extension config: %s\n", err)
	}
}

// GetExtensionByType 获取扩展
func GetExtensionByType(interfaceType reflect.Type, spiName string) (*plugin.Plugin, error) {
	theType := reflect2.Type2(interfaceType)
	i := allExtension[theType]
	if i != nil {
		extension := i[spiName]
		if extension != nil {
			return extension, nil
		}
	}
	return nil, errors.New(interfaceType.String() + " extension:" + spiName + " not find")
}

// GetExtensionByName 获取扩展
func GetExtensionByName(pkg string, structName string, spiName string) (*plugin.Plugin, error) {
	theType := reflect2.TypeByPackageName(pkg, structName)
	i := allExtension[theType]
	if i != nil {
		extension := i[spiName]
		if extension != nil {
			return extension, nil
		}
	}
	return nil, errors.New(pkg + " " + structName + " extension:" + spiName + " not find")
}

func Lookup(plugin *plugin.Plugin, methodName string) (plugin.Symbol, error) {
	if plugin == nil || len(methodName) <= 0 {
		return nil, errors.New("非法参数")
	}
	if cache[plugin] == nil {
		symbol, err := plugin.Lookup(methodName)
		if err != nil {
			return nil, err
		}
		cache[plugin] = symbol
	}
	return cache[plugin], nil
}

// LookupMethod 根据方法名称查找方法
func LookupMethod(interfaceType reflect.Type, spiName string, methodName string) (plugin.Symbol, error) {
	extension, _ := GetExtensionByType(interfaceType, spiName)
	method, err := Lookup(extension, methodName)
	if err != nil {
		return nil, err
	}
	if method == nil {
		return nil, errors.New("method:" + methodName + " can not find in:" + interfaceType.String() + " " + spiName)
	}
	return method, nil
}

// BuildDynamicLibrary 构建动态链接库
func BuildDynamicLibrary(targetPath string, sourcePath string) string {
	targetPath = GetDynamicLibraryName(targetPath)
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

func GetDynamicLibraryName(targetPath string) string {
	split := strings.Split(targetPath, ".")
	if runtime.GOOS == "windows" {
		targetPath = split[0] + ".dll"
	} else {
		targetPath = split[0] + ".so"
	}
	return targetPath
}
