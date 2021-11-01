package main

import "fmt"

func FindMinInArray(array []int) int {
	var minNumber int

	for index, current := range array {
		if index == 0 || current < minNumber {
			minNumber = current
		}
	}

	return minNumber
}

func main() {
	x := []int{48, 96, 86, 68, 57, 82, 63, 70, 37, 34, 83, 27, 19, 97, 9, 17}

	fmt.Println(fmt.Sprintf("Минимальное число %d", FindMinInArray(x)))
}
