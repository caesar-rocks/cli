package util

import (
	"fmt"
	"os"
)

func ExitWithError(err error) {
	fmt.Println(err)
	os.Exit(1)
}
