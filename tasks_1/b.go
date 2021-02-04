package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	var input string

	fmt.Scan(&input)

	count := utf8.RuneCountInString(input)
	summ := count * 23
	result := summ / 100

	if result <= 1 {
		fmt.Printf("%v коп.", summ%100)
	} else {
		fmt.Printf("%v р. %v коп.", result, summ%100)
	}
}
