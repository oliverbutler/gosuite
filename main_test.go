package main

import (
	"testing"
)

func sum(a, b int) int {
	return a + b
}

func TestSum(t *testing.T) {
	res := sum(1, 2)

	res2 := sum(2, 3)

	if res2 != 5 {
		t.Error("Expected 2 + 3 to equal 5")
	}

	if res != 3 {
		t.Error("Expected 1 + 2 to equal 3")
	}
}
