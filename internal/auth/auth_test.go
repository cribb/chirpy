// auth_test.go

package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
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

func TestJWT(t *testing.T) {
	secret := "WannaKnowASecret?" // "PromiseNotToTell"
	testUUID := uuid.New()
	fmt.Printf("TestUUID: %v\n", testUUID)

	jwtWithUUID, err := MakeJWT(testUUID, secret, 10*time.Second)
	if err != nil {
		fmt.Printf("Error generating authenticated UUID. (%v)", err)
	}
	fmt.Printf("jwt: %v\n", jwtWithUUID)

	authenticatedUUID, err := ValidateJWT(jwtWithUUID, secret)
	if err != nil {
		fmt.Printf("Error authenticating UUID. (%v)", err)
	}
	fmt.Printf("Authenticated UUID: %v\n", authenticatedUUID)

	if testUUID != authenticatedUUID {
		t.Errorf("authenication process for UUID has failed.")
	}
}
