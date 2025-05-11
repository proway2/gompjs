package main

import (
	"encoding/json"
	"fmt"

	"github.com/proway2/gomp/pkg/gompjs"
)

func main() {
	input := "{'chompjs': \"is a beautiful\", library: {created_by: 'Nykakin'}, 'this package': 'is just a wrapper', 'abc': 2025}"
	fmt.Printf("Input: %+v\n", input)
	fmt.Println("\nParsing with the standard library ...")
	var res any
	if err := json.Unmarshal([]byte(input), &res); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Result: %+v\n", res)
	}
	fmt.Println("\nParsing with `gompjs.ParseJsObject` ...")
	result, err := gompjs.ParseJsObject(&input, false, json.Unmarshal)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Result is of type: %T\n", result)
	fmt.Printf("Result: %+v\n", result)

	fmt.Println("\nParsing with `gompjs.ParseJsObjects` ...")
	complexInput := "{a: 'b', 'c': 21} ['d', 121, 'e']"
	fmt.Printf("\nComplex input: %+v\n", complexInput)
	data, _ := gompjs.ParseJsObjects(&complexInput, false, false, json.Unmarshal)
	for value := range data {
		fmt.Printf("Element1: %v, type: %T\n", value, value)
	}
}
