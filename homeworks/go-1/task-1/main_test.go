package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestMeteresToFeet(t *testing.T) {
	if MeteresToFeet(1) != 0.3048 {
		t.Error("ошибка, 1 метр равен 0.304810 футов")
	}

	if !strings.Contains(fmt.Sprintf("%.4f\n", MeteresToFeet(10)), "32.8084") {
		t.Error("ошибка, 10 метров равны 32.8084 футам")
	}
}
