package main

import (
	"fmt"
)

func IsDivThree(val int) bool {
	return val % 3 == 0
}

func main() {
	for index := range make([]int, 100) {
		if IsDivThree(index) {
			fmt.Println(index)
		}
	}
}

