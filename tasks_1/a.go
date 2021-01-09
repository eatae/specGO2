package main

import "fmt"

func main() {
	var weight int
	fmt.Scan(&weight)

	if weight != 2 && weight%2 == 0 {
		fmt.Print("YES")
	} else {
		fmt.Print("NO")
	}
}
