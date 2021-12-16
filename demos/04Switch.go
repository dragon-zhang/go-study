package demos

import "fmt"

func Switch(name int) {
	switch {
	case name < 0:
		fmt.Println("<0")
	case name == 0:
		fmt.Println("=0")
	case name > 0:
		fmt.Println(">0")
	}
}
