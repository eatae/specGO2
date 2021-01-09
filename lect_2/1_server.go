package main

import (
	"fmt"
	"log"
	"net/http"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from first page")
	fmt.Println("Homepage now linking")
}

func shutdown() {
	fmt.Println("Rest API shut down")
}

func main() {
	fmt.Println("Rest API V1 worked")
	defer shutdown()
	http.HandleFunc("/", HomePage)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
