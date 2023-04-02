package minifac

import (
	"fmt"
	"os"
)

func Log(pattern string, args ...any) {
	fmt.Fprintf(os.Stdout, pattern+"\n", args...)
}
