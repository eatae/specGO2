package main

import "fmt"

func main() {
	var (
		one int
		two int
	)
	fmt.Scan(&one)
	fmt.Scan(&two)

	fmt.Println((one + two) * (one + two))
}
