// auth_test.go

package auth

import (
	"fmt"
	"testing"

	"github.com/alexedwards/argon2id"
)

func TestCheckHashPassword(t *testing.T) {
	password := "12345678"
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	fmt.Printf("plaintext password: %v\n", password)
	if err != nil {
		fmt.Printf("Error Creating Hash: %v", err)
	}
	if err == nil {
		fmt.Printf("hashed password: %v\n", hash)
	}

	match, err := argon2id.ComparePasswordAndHash(password, hash)

	if !match || err != nil {
		t.Errorf("password has and hashed(password) do not match.")
	}
}
