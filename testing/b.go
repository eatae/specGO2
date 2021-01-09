package main

import "fmt"

func main() {
	var (
		one int
		two int
	)
	fmt.Scan(&one)
	fmt.Scan(&two)

	fmt.Printf("Периметр: %d \n", (one+two)*2)
	fmt.Printf("Площадь: %d \n", one*two)
}
