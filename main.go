package main

import (
	"os"

	"github.com/oliversabler/apa/repl"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
