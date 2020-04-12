package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/jakubbartel/kbcdatastudioproc"
)

func main() {
	if err := kbcdatastudioproc.RunE(); errors.Is(err, kbcdatastudioproc.ErrUser) {
		_, _ = fmt.Fprintf(os.Stderr, "User error: %v", err)
		os.Exit(1)
	} else if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Application error: %v", err)
		os.Exit(2)
	}

	os.Exit(0)
}
