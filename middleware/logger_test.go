package middleware

import (
	"log"
	"testing"
)

func TestMakeSensitive(t *testing.T) {
	pass1 := `{"key":"1","password11" 	: "pass1","a":{"password":"pass212  "}}`
	log.Println(makeSensitive(pass1))
}
