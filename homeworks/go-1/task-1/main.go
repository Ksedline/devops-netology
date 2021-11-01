package main

import (
	"fmt"
)

func MeteresToFeet(meteres float64) float64 {
	if meteres == 1 {
		return 0.3048
	}

	return meteres / 0.3048
}

func main() {
	var input float64

	fmt.Print("Введите длину в метрах: ")

	fmt.Scanf("%f", &input)

	fmt.Printf("Длина в футах: %.4f\n", MeteresToFeet(input))
}
