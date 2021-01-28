package main

import (
	"fmt"
	"os"

	"github.com/jameslahm/glox"
)

func main() {
	args := os.Args
	var g = &glox.Glox{}
	if len(args) > 2 {
		fmt.Println("Usage: glox [script]")
	} else if len(args) == 2 {
		g.RunFile(args[1])
	} else {
		g.RunPrompt()
	}
}
