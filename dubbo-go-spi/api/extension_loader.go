package api

import (
	"bytes"
	"errors"
	"github.com/modern-go/reflect2"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"
	"runtime"
	"strings"
)

type ExtensionLoader struct {
	packageName      string
	structName       string
	uniqueName       string
	cachedExtensions map[string]*plugin.Plugin
	cached           map[string]interface{}
	dirs             []string
}

var splitter = "@"
var dubboDirectory = "META-INF/dubbo/"
var dubboInternalDirectory = dubboDirectory + "internal/"
var extensionLoaders = make(map[string]*ExtensionLoader)

//todo 校验packageName、structName的正确性
func getUniqueName(packageName string, structName string) string {
	if EndsWith(packageName, "/") {
		return packageName[0:len(packageName)-1] + splitter + structName
	}
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

// GetExtensionLoader 获取扩展加载器
func GetExtensionLoader(packageName string, structName string, dirs ...string) (*ExtensionLoader, error) {
	if len(packageName) <= 0 || len(structName) <= 0 {
		return nil, errors.New("packageName or structName is empty")
	}
	uniqueName := getUniqueName(packageName, structName)
	if extensionLoaders[uniqueName] == nil {
		extensionLoaders[uniqueName] = &ExtensionLoader{
			packageName: packageName, structName: structName, uniqueName: uniqueName,
			cachedExtensions: make(map[string]*plugin.Plugin),
			cached:           make(map[string]interface{}), dirs: dirs}
	}
	return extensionLoaders[uniqueName], nil
}

// GetExtensions 获取所有扩展
func (loader *ExtensionLoader) GetExtensions() map[string]*plugin.Plugin {
	if len(loader.cachedExtensions) == 0 {
		loader.cachedExtensions = loader.LoadExtensions()
	}
	return loader.cachedExtensions
}

// LoadExtensions 加载扩展
func (loader *ExtensionLoader) LoadExtensions() map[string]*plugin.Plugin {
	var extensions = make(map[string]*plugin.Plugin)
	loader.LoadDirectory(extensions, dubboDirectory)
	loader.LoadDirectory(extensions, dubboInternalDirectory)
	for _, dir := range loader.dirs {
		loader.LoadDirectory(extensions, dir)
	}
	return extensions
}

// LoadDirectory 读取资源目录
func (loader *ExtensionLoader) LoadDirectory(extensions map[string]*plugin.Plugin, dir string) {
	_ = filepath.Walk(dir,
		func(path string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if f.IsDir() {
				return nil
			}
			if strings.ReplaceAll(loader.uniqueName, "/", "#") != path[len(dir+"/"):] {
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
			pkg := strings.ReplaceAll(split[0][len(dir+"/"):], "#", "/")
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
		extension := split[1]
		if EndsWith(extension, ".go") {
			//需要编译
			plug, err := loader.AddExtension(spiName, extension)
			if err == nil {
				extensions[spiName] = plug
				continue
			}
		} else {
			plug, err := loader.AddExtension(spiName, extension+suffix[0])
			if err == nil {
				extensions[spiName] = plug
				continue
			}
			plug, err = loader.AddExtension(spiName, extension+suffix[1])
			if err == nil {
				extensions[spiName] = plug
				continue
			}
			return err
		}
	}
	return nil
}

// AddExtension 添加扩展
func (loader *ExtensionLoader) AddExtension(name string, extension string) (*plugin.Plugin, error) {
	suffix, err := GetDylibSuffix(extension)
	if err != nil {
		return nil, err
	}
	if suffix == "go" {
		dylib, err := BuildDylib(extension)
		if err != nil {
			return nil, err
		}
		return loader.CreateExtension(name, dylib)
	}
	suffixes, err := GetCorrectDylibSuffix()
	if err != nil {
		return nil, err
	}
	index := Contains(suffixes, suffix)
	if index == -1 {
		return nil, errors.New("Illegal dynamic link library suffix:" + suffix)
	}
	return loader.CreateExtension(name, extension)
}

// CreateExtension 加载动态链接库
func (loader *ExtensionLoader) CreateExtension(name string, dylib string) (*plugin.Plugin, error) {
	result := loader.cached[name]
	if result == nil {
		plug, err := plugin.Open(dylib)
		if err == nil {
			loader.cached[name] = plug
		} else {
			loader.cached[name] = err
		}
		result = loader.cached[name]
	}
	return getResult(result)
}

func (loader *ExtensionLoader) GetLoadedExtensions() []string {
	keys := make([]string, 0, len(loader.cachedExtensions))
	for k := range loader.cachedExtensions {
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

func EndsWith(source string, target string) bool {
	return strings.LastIndex(source, target) == len(source)-len(target)
}

func Contains(array []string, target string) int {
	for i, s := range array {
		if s == target {
			return i
		}
	}
	return -1
}

// BuildDylib 延迟构建动态链接库
func BuildDylib(sourcePath string) (string, error) {
	split := strings.Split(sourcePath, ".")
	split = strings.Split(split[0], "/")
	suffix, _ := GetCorrectDylibSuffix()
	targetPath := dubboDirectory + split[len(split)-1] + suffix[0]
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
		return "", err
	}
	return targetPath, nil
}
