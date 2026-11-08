package auth

import (
	"fmt"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", nil
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

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	now := time.Now().UTC()
	expiresAt := now.Add(expiresIn)

	claims := &jwt.RegisteredClaims{
		Subject:   userID.String(),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		Issuer:    "chirpy-access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		fmt.Printf("Error occurred while signing token: %v\n", err)
		return "", err
	}
	fmt.Printf("tokenString: %v\n", tokenString)

	return tokenString, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

	var rclaims jwt.RegisteredClaims

	tokfunc := func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &rclaims, tokfunc)
	if err != nil {
		fmt.Printf("Error Parsing Secret. (%v)\n", err)
		return uuid.Nil, err
	}

	userid, err := token.Claims.GetSubject()
	if err != nil {
		fmt.Printf("Error handling claims. (%v)\n", err)
		return uuid.Nil, err
	}

	userUUID, err := uuid.Parse(userid)
	if err != nil {
		fmt.Printf("Error parsing UUID. (%v)\n", err)
		return uuid.Nil, err
	}
	return userUUID, nil
}
