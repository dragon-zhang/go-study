package demos

import "fmt"

func For() {
	for i := 0; i < 5; i++ {
		fmt.Println(i)
	}

	var m = map[string]int{
		"go":     1,
		"java":   2,
		"kotlin": 3,
	}

	for k, v := range m {
		fmt.Printf("key:%v, value:%v\n", k, v)
	}
}
