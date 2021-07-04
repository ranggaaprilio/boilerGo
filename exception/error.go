package exception

import (
	"fmt"
)

//PanicIfNeeded for handle error
func PanicIfNeeded(err interface{}) {
	defer recover()
	if err != nil {
		panic(err)

	}
}

func Catch() {
	if r := recover(); r != nil {
		fmt.Println("Error occured", r)
	}
}
