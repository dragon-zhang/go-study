package demos

import (
	"fmt"
)

func IfElse() {
	const name = 1
	if name < 0 {
		fmt.Println("<0")
	} else if name == 0 {
		fmt.Println("=0")
	} else {
		fmt.Println(">0")
	}
}
