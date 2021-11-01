package main

import "testing"

func TestCheckDivThree(t *testing.T) {
	if IsDivThree(2) {
		t.Error("ошибка, число 2 не делится на 3 без остатка")
	}
}
