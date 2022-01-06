package api

import (
	"errors"
	"github.com/modern-go/reflect2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"plugin"
	"runtime"
	"strings"
	"sync"
)

type ExtensionLoader struct {
	packageName string
	structName  string
	uniqueName  string
}

var splitter = "@"
var dubboDirectory = "META-INF/dubbo/"
var dubboInternalDirectory = dubboDirectory + "internal/"

//todo 校验packageName、structName的正确性
func getUniqueName(packageName string, structName string) string {
	return packageName + splitter + structName
}

func getResult(result interface{}) (*plugin.Plugin, error) {
	switch result.(type) {
	case error:
		return nil, result.(error)
	case *plugin.Plugin:
		return result.(*plugin.Plugin), nil
	default:
		return nil, errors.New("unknown error ")
	}
}

var extensionLoaders = make(map[string]*ExtensionLoader)

// GetExtensionLoader 获取扩展加载器
func GetExtensionLoader(packageName string, structName string) (*ExtensionLoader, error) {
	if len(packageName) <= 0 || len(structName) <= 0 {
		return nil, errors.New("packageName or structName is empty")
	}
	uniqueName := getUniqueName(packageName, structName)
	if extensionLoaders[uniqueName] == nil {
		extensionLoaders[uniqueName] = &ExtensionLoader{
			packageName: packageName, structName: structName, uniqueName: uniqueName}
	}
	return extensionLoaders[uniqueName], nil
}

var cachedExtensions = make(map[string]*plugin.Plugin)
var cached = make(map[string]interface{})

// GetExtensions 获取所有扩展
func (loader *ExtensionLoader) GetExtensions() map[string]*plugin.Plugin {
	if len(cachedExtensions) == 0 {
		cachedExtensions = loader.LoadExtensions()
	}
	return cachedExtensions
}

// LoadExtensions 加载扩展
func (loader *ExtensionLoader) LoadExtensions() map[string]*plugin.Plugin {
	var extensions = make(map[string]*plugin.Plugin)
	loader.LoadDirectory(extensions, dubboDirectory)
	loader.LoadDirectory(extensions, dubboInternalDirectory)
	return extensions
}

// LoadDirectory 读取资源目录
func (loader *ExtensionLoader) LoadDirectory(extensions map[string]*plugin.Plugin, dir string) {
	err := filepath.Walk(dir,
		func(path string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if f.IsDir() {
				return nil
			}
			if strings.ReplaceAll(loader.uniqueName, "/", "#") != path {
				//不是此Loader的职责
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
			return loader.LoadResource(extensions, path)
		})
	if err != nil {
		log.Fatalf("Error read extension config: %s\n", err)
	}
}

// LoadResource 加载配置文件，解析配置
func (loader *ExtensionLoader) LoadResource(extensions map[string]*plugin.Plugin, path string) error {
	suffix, err := GetCorrectDylibSuffix()
	if err != nil {
		return err
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	for _, line := range strings.Split(string(data), "\n") {
		split := strings.Split(line, "=")
		spiName := split[0]
		plug, err := loader.AddExtension(spiName, split[1]+suffix[0])
		if err == nil {
			extensions[spiName] = plug
			continue
		}
		plug, err = loader.AddExtension(spiName, split[1]+suffix[1])
		if err == nil {
			extensions[spiName] = plug
			continue
		}
		return err
	}
	return nil
}

// AddExtension 添加扩展
func (loader *ExtensionLoader) AddExtension(name string, dylib string) (*plugin.Plugin, error) {
	suffix, err := GetDylibSuffix(dylib)
	if err != nil {
		return nil, err
	}
	suffixes, err := GetCorrectDylibSuffix()
	if err != nil {
		return nil, err
	}
	index := Contains(suffixes, suffix)
	if index == -1 {
		return nil, errors.New("Illegal dynamic link library suffix:" + suffix)
	}
	return loader.CreateExtension(name, dylib)
}

var once = sync.Once{}

// CreateExtension 如果必要，创建扩展
func (loader *ExtensionLoader) CreateExtension(name string, dylib string) (*plugin.Plugin, error) {
	//只加载一次动态链接库就可以了
	once.Do(func() {
		plug, err := plugin.Open(dylib)
		if err == nil {
			cached[name] = plug
		} else {
			cached[name] = err
		}
	})
	result := cached[name]
	return getResult(result)
}

func (loader *ExtensionLoader) GetLoadedExtensions() []string {
	keys := make([]string, 0, len(cachedExtensions))
	for k := range cachedExtensions {
		keys = append(keys, k)
	}
	return keys
}

func (loader *ExtensionLoader) GetSupportedExtensions() []string {
	extensions := loader.GetExtensions()
	keys := make([]string, 0, len(extensions))
	for k := range extensions {
		keys = append(keys, k)
	}
	return keys
}

func (loader *ExtensionLoader) GetExtension(name string) *plugin.Plugin {
	return loader.GetExtensions()[name]
}

func GetDylibSuffix(dylib string) (string, error) {
	if len(dylib) <= 0 {
		return "", errors.New("Illegal dynamic link library name ")
	}
	split := strings.Split(dylib, ".")
	if len(split) <= 1 {
		return "", errors.New("Dynamic link library must have a suffix ")
	}
	return split[0], nil
}

func GetCorrectDylibSuffix() ([]string, error) {
	theOs := runtime.GOOS
	if theOs == "darwin" {
		return []string{".dylib", ".so"}, nil
	} else if theOs == "windows" {
		return []string{".dll"}, nil
	} else if theOs == "linux" {
		return []string{".so"}, nil
	}
	return []string{}, errors.New("unsupported os ")
}

func Contains(array []string, target string) int {
	for i, s := range array {
		if s == target {
			return i
		}
	}
	return -1
}
