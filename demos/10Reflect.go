package demos

import (
	"fmt"
	"reflect"
)

// Reflect interface{}可以接手任何值，类似java的Object
func Reflect(args ...interface{}) {
	for _, v := range args {
		fmt.Print(reflect.TypeOf(v))
		fmt.Print(" ")
	}
	fmt.Println()
}
