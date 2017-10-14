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
