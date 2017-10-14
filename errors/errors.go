package errors

import (
	"fmt"
	"os"
)

func ErrCheck(err error) {
	if err != nil {
		fmt.Println("err: ", err.Error())
		os.Exit(-1)
	}
}

func NilCheck(v interface{}) {
	if v == nil {
		fmt.Printf("%#v is nil", v)
		os.Exit(-1)
	}
}
