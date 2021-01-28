package utils

import "fmt"

func Error(line int, message string) {
	fmt.Printf("[line %d]  Error %s\n", line, message)
}
