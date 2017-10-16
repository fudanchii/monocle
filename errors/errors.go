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

func AssertFalse(cond bool, msg string) {
	if cond {
		fmt.Println("err: ", msg)
		os.Exit(-1)
	}
}

func AssertTrue(cond bool, msg string) {
	if !cond {
		fmt.Println("err: ", msg)
		os.Exit(-1)
	}
}
