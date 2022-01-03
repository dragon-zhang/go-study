package demos

import (
	"fmt"
	"reflect"
)

type TagTest struct {
	f1 string "f1:\"one\""
	f2 string `f2:"three"`
	f3 string `one:"1" two:"2" blank:""`
}

func Tag() {
	t := reflect.TypeOf(TagTest{})
	f1, _ := t.FieldByName("f1")
	fmt.Println(f1.Tag)
	f2, _ := t.FieldByName("f2")
	fmt.Println(f2.Tag)

	f3, _ := t.FieldByName("f3")
	fmt.Println(f3.Tag) // one:"1" two:"2"blank:""
	value, find := f3.Tag.Lookup("one")
	fmt.Printf("%s, %t\n", value, find) // 1, true
	value, find = f3.Tag.Lookup("blank")
	fmt.Printf("%s, %t\n", value, find) // , true
	value, find = f3.Tag.Lookup("five")
	fmt.Printf("%s, %t\n", value, find) // , false
}
