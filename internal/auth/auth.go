package auth

import (
	"fmt"
	"log"

	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {
	fmt.Printf("Ho dee ho ho.\n")

	var params *argon2id.Params

	params = argon2id.DefaultParams

	hash, err := argon2id.CreateHash(password, params)
	if err != nil {
		log.Fatal(err)
	}
	return hash, nil

}

func CheckPasswordHash(password string, hash string) (bool, error) {
	chk, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return chk, nil
}
