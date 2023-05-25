package utils

import (
	"fmt"
	"math/rand"
	"strings"
)

const alpha = "abcdefghijklmnopqrstuvwxyz"
const alphanum = "abcdefghijklmnopqrstuvwxyz1234567890"

func RandomAlphabet(n int) string {
	var sb strings.Builder
	l := len(alpha)

	for i := 0; i < n; i++ {
		c := alpha[rand.Intn(l)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomAlphanum(n int) string {
	var sb strings.Builder
	l := len(alphanum)

	for i := 0; i < n; i++ {
		c := alphanum[rand.Intn(l)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomUsername() string {
	return RandomAlphabet(8)
}

func RandomPassword() string {
	return RandomAlphanum(10)
}

func RandomEmail() string {
	username := RandomUsername()
	return fmt.Sprintf("%v@email.com", username)
}

func RandomFullname() string {
	firstName := RandomAlphabet(6)
	surname := RandomAlphabet(6)
	return fmt.Sprintf("%v %v", firstName, surname)
}
