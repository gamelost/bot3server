package main

import (
	"fmt"
	"testing"
)

func TestSanity(t *testing.T) {
	const a, b = 2, 2
	if a+b == 4 {
		fmt.Println("All working")
	}
}

func TestAddHandler(t *testing.T) {
}
