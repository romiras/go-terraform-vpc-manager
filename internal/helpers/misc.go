package helpers

import (
	"fmt"
	"os"
)

func DebugMsg(m ...interface{}) {
	fmt.Println(m...)
}

func AbortOnError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
