package main

import (
	"testing"
)

func TestMinFromArray(t *testing.T) {
	val := []int{10, 8, 20}

	minNumber := FindMinInArray(val)

	if minNumber != 8 {
		t.Error("ошибка, минимальное число найдено неправильно")
	}
}
