package eprintf

import (
	"fmt"
	"os"
)

/**
 * Printf to stderr and exit with -1
 */
func EPrintf(format string, a ...any) {
	fmt.Fprintln(os.Stderr, fmt.Errorf(format, a...))
	os.Exit(-1)
}
