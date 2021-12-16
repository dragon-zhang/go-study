package demos

//Function
/*
 跟java一样，变参只能放在最后
 注意默认是值用传递，而不是引用传递
*/
func Function(arg ...int) (int, bool) {
	var result = 0
	for _, v := range arg {
		result += v
	}
	return result, true
}
