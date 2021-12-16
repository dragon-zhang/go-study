package demos

import (
	"fmt"
	json "github.com/json-iterator/go"
)

type User struct {
	Account  string
	Password string
}

func ToJson(obj interface{}) string {
	bytes, err := json.ConfigCompatibleWithStandardLibrary.Marshal(obj)
	if err != nil {
		return "{}"
	}
	return string(bytes)
}

func FromJson(s string, result interface{}) interface{} {
	err := json.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(s), &result)
	if err != nil {
		return nil
	}
	return result
}

func Json() {
	user := &User{
		Account:  "admin",
		Password: "123456",
	}
	toJson := ToJson(user)
	fmt.Println(toJson)

	result := &User{}
	FromJson(toJson, result)
	fmt.Println(result.Account)
	fmt.Println(result.Password)
}
