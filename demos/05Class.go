package demos

import "fmt"

type Class struct {
	name string
	age  int
}

func (c *Class) GetName() string {
	return c.name
}

func (c *Class) GetAge() int {
	return c.age
}

// SetName 支持链式调用
func (c *Class) SetName(name string) *Class {
	c.name = name
	return c
}

// SetAge 支持链式调用
func (c *Class) SetAge(age int) *Class {
	c.age = age
	return c
}

func newClass() Class {
	return Class{}
}

func TestNewClass() {
	class := newClass()
	class.SetName("zhangsan").
		SetAge(10)
	fmt.Println(class.GetName())
	fmt.Println(class.GetAge())
}
