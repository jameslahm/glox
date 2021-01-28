package utils

import "fmt"

const UNEXPECTED_CHARACTER_MESSAGE = "Unexpected character"
const UNTERMINATED_STRING = "Unterminated string"
const INVALID_NUMBER = "Invalid number"
const UNMATCHED_PAREN = "Unmatched paren"

func Error(line int, message string) {
	fmt.Printf("[line %d]  Error %s", line, message)
}
