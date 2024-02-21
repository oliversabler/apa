package main

import (
	"os"

	"github.com/oliversabler/egglang/repl"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
