package authlib

import (
	"github.com/barreyo/esportsdrafts/libs/log"
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
	_ = JWTMiddleware(
		JWTConfig{
			SigningKey:  []byte("asdf"),
			AllowedRole: "user",
		},
	)
	// TODO
}

func TestGenTestTokens(t *testing.T) {
	claims := &JWTClaims{
		Username: "pelle",
		UserID:   "random_id",
		Roles: []string{
			"user",
		},
	}
	token, _, err := GenerateAuthToken(claims, 60*time.Minute, []byte("secret"))
	if err != nil {
		t.Errorf("Something went wrong generating token. Error: %e", err)
	}
	log.GetLogger().Infof("T1: %s", token)

	claims.Username = "race"
	token, _, err = GenerateAuthToken(claims, 60*time.Minute, []byte("secret"))
	if err != nil {
		t.Errorf("Something went wrong generating token. Error: %e", err)
	}
	log.GetLogger().Infof("T2: %s", token)
}

func TestJWTRace(t *testing.T) {
	e := echo.New()
	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}
	initialToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InBlbGxlIiwidXNlcl9pZCI6InJhbmRvbV9pZCIsInJvbGVzIjpbInVzZXIiXSwiZXhwIjoxNTczNzUyNzkyLCJqdGkiOiI3NjIyNTY3OS1jMTVkLTQ3OTAtYmE5ZC00NWNkOTJmZmRmZDUiLCJpYXQiOjE1NzM3NDkxOTJ9.JEU53TTA3yNAJGc0pjrD_s2nwQrPGcZabEkKCKTo1dk"
	raceToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InJhY2UiLCJ1c2VyX2lkIjoicmFuZG9tX2lkIiwicm9sZXMiOlsidXNlciJdLCJleHAiOjE1NzM3NTI3OTIsImp0aSI6ImVkNjY2MGRlLWQxNTQtNGYxNy05NDU3LTA2NDQyZDc0Y2NlMiIsImlhdCI6MTU3Mzc0OTE5Mn0.oYv-AD4n0e42YgdRvfYggTQua7wCnjwfRlKHufpXXJY"
	validKey := []byte("secret")

	h := JWTMiddleware(JWTConfig{
		AllowedRole: "user",
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

	makeReq(initialToken)
	// user := c.Get("user").(*jwt.Token)
	// _ms := user.Claims.(*JWTClaims)
	// assert.Equal(t, claims.Name, "John Doe")

	makeReq(raceToken)
	// user = c.Get("user").(*jwt.Token)
	// claims = user.Claims.(*jwtCustomClaims)
	// Initial context should still be "John Doe", not "Race Condition"
	// assert.Equal(t, claims.Name, "John Doe")
	// assert.Equal(t, claims.Admin, true)
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
