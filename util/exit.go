package util

import (
	"fmt"
	"os"
)

func ExitWithError(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func ExitAndCleanUp(path string, err error) {
	if removeErr := os.RemoveAll(path); removeErr != nil {
		ExitWithError(
			fmt.Errorf(
				"error while cleaning up: %s, and error while setting up: %s",
				removeErr,
				err,
			),
		)
	} else {
		ExitWithError(err)
	}
}
