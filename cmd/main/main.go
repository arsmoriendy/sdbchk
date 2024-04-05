package main

import (
	"github.com/arsmoriendy/sdbchk"
	"os"
)

func main() {
	sdbchk.SdbChk(os.Args[1], os.Args[2])
}
