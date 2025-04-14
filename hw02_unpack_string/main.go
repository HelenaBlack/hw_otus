package main

import "fmt"

func main() {
	var input string
	fmt.Print("Enter a string: ")
	_, err := fmt.Scanln(&input)
	if err != nil {
		fmt.Println(err)
		return
	}
	result, err := Unpack(input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)
}
