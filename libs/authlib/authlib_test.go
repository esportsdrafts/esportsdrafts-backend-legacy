package authlib

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
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

func TestJWTRace(t *testing.T) {
	e := echo.New()
	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}
	initialToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"
	raceToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IlJhY2UgQ29uZGl0aW9uIiwiYWRtaW4iOmZhbHNlfQ.Xzkx9mcgGqYMTkuxSCbJ67lsDyk5J2aB7hu65cEE-Ss"
	validKey := []byte("secret")

	h := JWTMiddleware(JWTConfig{
		AllowedRole: "pelle",
		SigningKey:  validKey,
	})(handler)

	makeReq := func(token string) echo.Context {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()
		req.Header.Set(echo.HeaderAuthorization, "Bearer: "+token)
		c := e.NewContext(req, res)

		reqResult := h(c)
		if reqResult != nil {
			t.Errorf("Failed to perfom request: %e", reqResult)
		}
		return c
	}

	c := makeReq(initialToken)
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JWTClaims)
	assert.Equal(t, claims.Name, "John Doe")

	makeReq(raceToken)
	user = c.Get("user").(*jwt.Token)
	claims = user.Claims.(*jwtCustomClaims)
	// Initial context should still be "John Doe", not "Race Condition"
	assert.Equal(t, claims.Name, "John Doe")
	assert.Equal(t, claims.Admin, true)
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
