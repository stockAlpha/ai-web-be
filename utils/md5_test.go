package utils

import (
	"fmt"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password, _ := HashPassword("")
	fmt.Println("password", password)
}
