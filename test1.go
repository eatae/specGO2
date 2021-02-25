package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func main() {
	fileData, err := ioutil.ReadFile("users.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	var parsedData map[string]interface{}
	json.Unmarshal(fileData, &parsedData)

	for key, value := range parsedData {
		fmt.Println("key:", key, "value:", value)
	}
}
