// auth_test.go

package auth

import (
	"fmt"
	"net/http"
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

func TestGetBearerToken(t *testing.T) {

	token := "SecretsAndLies"

	h := http.Header{}
	_, err := GetBearerToken(h)
	if err == nil {
		t.Errorf("Test for empty header failed")
	}

	h.Set("Authorization", "Bearer "+token)
	testToken, err := GetBearerToken(h)
	if err != nil || testToken != token {
		fmt.Printf("Something weird happened with the Bearer token. (%v)\n", err)
		fmt.Printf("token: %v, testToken: %v \n", token, testToken)
		t.Fatalf("Bearer token test failed.")
	}

}

func TestMakeRefreshToken(t *testing.T) {
	// stupid test just to see if the keygen is working
	Nbytes := 32 // must be even
	tokHexString := MakeRefreshToken(Nbytes)

	if Nbytes != len(tokHexString)/2 {
		t.Errorf("Test for HexString byte-length failed. (%v)", tokHexString)
	}
}
