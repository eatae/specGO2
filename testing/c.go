package main

import "fmt"

func main() {
	var (
		one   int
		two   int
		three int
	)
	fmt.Scan(&one)
	fmt.Scan(&two)
	fmt.Scan(&three)

	fmt.Printf("%d%d%d\n", three, two, one)
}
