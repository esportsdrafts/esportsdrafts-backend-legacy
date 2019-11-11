package authlib

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func TestReconstructToken(t *testing.T) {
	words := []struct {
		input       string
		input2      string
		shouldEqual string
	}{
		{"pelle", "boi", "pelle.boi"},
		{"dance", "game", "dance.game"},
		{"DanCe", "gamE", "DanCe.gamE"},
	}
	for _, table := range words {
		res := reconstructAuthToken(table.input, table.input2)
		if res != table.shouldEqual {
			t.Errorf("Checking token reconstruct ('%s', '%s') was incorrect, got %s, wanted %s", table.input, table.input2, res, table.shouldEqual)
		}
	}
}

func TestCreateMiddleware(t *testing.T) {
	middlewareFunc := JWTMiddleware(
		JWTConfig{
			SigningKey:  []byte("asdf"),
			AllowedRole: "user",
		},
	)
	// TODO
}

func TestShouldSetCookie(t *testing.T) {
	// TODO
}

func TestGenerateAuthToken(t *testing.T) {
	claims := &JWTClaims{
		Username: "pelle",
		UserID:   "random_id",
		Roles: []string{
			"user",
		},
	}
	token, expiry, err := GenerateAuthToken(claims, 60*time.Minute, []byte("fake_key"))
	if err != nil {
		t.Errorf("Something went wrong generating token. Error: %e", err)
	}

	readToken, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("fake_key"), nil
	})

	readClaims, ok := readToken.Claims.(*JWTClaims)

	if !ok {
		t.Error("Read token not ok")
	}

	if err != nil {
		t.Errorf("Failed to parse token. Error: %e", err)
	}

	if !readToken.Valid {
		t.Error("Reading token failed")
	}

	if !contains(readClaims.Roles, "user") {
		t.Errorf("Wrong claims. Got: %s, wanted: [\"user\"]", readClaims.Roles)
	}

	expTimeUnix := expiry.Unix()
	readExpTime := readClaims.StandardClaims.ExpiresAt
	if expTimeUnix != readExpTime {
		t.Errorf("Expiry in claims different from plain token expiration. Got: %d, wanted: %d", expTimeUnix, readExpTime)
	}
}
