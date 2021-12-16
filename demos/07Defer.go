package demos

import "fmt"

var Name = "go"

// Defer defer类似java里的finally
func Defer() string {
	//先输出2，再输出1
	defer printf("1")
	defer printf("2")
	defer func() {
		Name = "python"
		fmt.Println("defer->", Name)
	}()
	fmt.Println("test->", Name)
	return Name
}

func printf(args ...string) {
	fmt.Println(args)
}
