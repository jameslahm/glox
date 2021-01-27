package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	args := os.Args
	if len(args) > 2 {
		fmt.Println("Usage: glox [script]")
	} else if len(args) == 2 {
		runFile(args[1])
	} else {
		runPrompt()
	}
}

func runFile(path string) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("ReadFile error: %v\n", err)
	}
	run(string(buf))
}

func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(">> ")
	for scanner.Scan() {
		line := scanner.Text()
		run(line)
		fmt.Print(">> ")
	}
}

func run(script string) {
	scanner := bufio.NewScanner(bytes.NewBuffer([]byte(script)))
	for scanner.Scan() {
		token := scanner.Text()
		fmt.Printf("%s ", token)
	}
	fmt.Println()
}
