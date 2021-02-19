package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

// action
func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Welcome to Home Page RESTAPI</h1>")
	dump, _ := httputil.DumpRequest(r, true)
	fmt.Fprintf(w, "<pre>")
	fmt.Fprintf(w, "%q", dump)
	fmt.Fprintf(w, "</pre>")
	fmt.Println("Someone visited the HomePage...")
}

func EndServe() {
	fmt.Println("---Stop server----")
}

func main() {
	defer EndServe()
	fmt.Println("----Start server----")
	// set route
	http.HandleFunc("/", HomePage)
	// start serve
	log.Fatal(http.ListenAndServe(":8050", nil))
	fmt.Println("----Stop----")
}
