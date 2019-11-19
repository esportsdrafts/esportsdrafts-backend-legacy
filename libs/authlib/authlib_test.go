package authlib

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/esportsdrafts/esportsdrafts/libs/log"

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

func TestPrintTestTokens(t *testing.T) {
	claims := &JWTClaims{
		Username: "pelle",
		UserID:   "random_id",
		Roles: []string{
			"user",
		},
	}
	// 5 years
	token, _, err := GenerateAuthToken(claims, 5*8760*time.Hour, []byte("secret"))
	if err != nil {
		t.Errorf("Something went wrong generating token. Error: %e", err)
	}
	log.GetLogger().Infof("T1: %s", token)

	claims.Username = "race"
	token, _, err = GenerateAuthToken(claims, 5*8760*time.Minute, []byte("secret"))
	if err != nil {
		t.Errorf("Something went wrong generating token. Error: %e", err)
	}
	log.GetLogger().Infof("T2: %s", token)

	claims.Username = "invalid-token"
	token, _, err = GenerateAuthToken(claims, 1*time.Second, []byte("secret"))
	if err != nil {
		t.Errorf("Something went wrong generating token. Error: %e", err)
	}
	log.GetLogger().Infof("T3: %s", token)
}

func TestShouldPanicWithoutSigningSecret(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Did not panic when no siging secret provided")
		}
	}()

	// Without config at all
	JWTMiddleware(JWTConfig{})

	// Explicit nil
	JWTMiddleware(JWTConfig{
		SigningKey: nil,
	})
}

func TestJWTRace(t *testing.T) {
	e := echo.New()
	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}
	initialToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InBlbGxlIiwidXNlcl9pZCI6InJhbmRvbV9pZCIsInJvbGVzIjpbInVzZXIiXSwiZXhwIjoxNzMxNDMzMzA0LCJqdGkiOiIwOTNjMTY5Yi04Y2U0LTRmOTctYWY3Ni1mMWIwMGE5YzdhYzciLCJpYXQiOjE1NzM3NTMzMDR9.PMQLCMJIiXQhPDl9cszRPo0bSqnB5e7-3OvoKNLBKLc"
	secondToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InJhY2UiLCJ1c2VyX2lkIjoicmFuZG9tX2lkIiwicm9sZXMiOlsidXNlciJdLCJleHAiOjE1NzYzODMzMDIsImp0aSI6ImZlZWUyZGMxLWRmMTAtNDc4YS04NTI0LWM5YzRkOWNlN2JjMiIsImlhdCI6MTU3Mzc1NTMwMn0.1wLsU99h0GqIcv_Nxip3Ede3_fcsOlJJ1xutZxDO65g"
	validKey := []byte("secret")

	h := JWTMiddleware(JWTConfig{
		SigningKey: validKey,
	})(handler)

	makeReq := func(token string) echo.Context {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()
		req.Header.Set(echo.HeaderAuthorization, "Bearer: "+token)
		c := e.NewContext(req, res)

		reqResult := h(c)
		if reqResult != nil {
			t.Errorf("Failed to perfom request: %+v", reqResult)
		}
		return c
	}

	c := makeReq(initialToken)
	user := c.Get("user").(*JWTClaims)
	if user.Username != "pelle" {
		t.Errorf("Username in initial token not matching. Wanted %s, got %s", "pelle", user.Username)
	}

	c = makeReq(secondToken)
	user = c.Get("user").(*JWTClaims)
	// Initial context should still be "John Doe", not "Race Condition"
	if user.Username != "race" {
		t.Errorf("Username should be 'race', got '%s'", user.Username)
	}
}

func TestShouldSetCookie(t *testing.T) {
	// TODO
}

func TestNoRoleMatch(t *testing.T) {
	e := echo.New()
	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}
	initialToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InBlbGxlIiwidXNlcl9pZCI6InJhbmRvbV9pZCIsInJvbGVzIjpbInVzZXIiXSwiZXhwIjoxNzMxNDMzMzA0LCJqdGkiOiIwOTNjMTY5Yi04Y2U0LTRmOTctYWY3Ni1mMWIwMGE5YzdhYzciLCJpYXQiOjE1NzM3NTMzMDR9.PMQLCMJIiXQhPDl9cszRPo0bSqnB5e7-3OvoKNLBKLc"
	secondToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InJhY2UiLCJ1c2VyX2lkIjoicmFuZG9tX2lkIiwicm9sZXMiOlsidXNlciJdLCJleHAiOjE1NzM3NTY5MDQsImp0aSI6ImI3ZDg0MjhmLWM3MWItNDk3Zi1iNzI2LThlYjBiYjQzODlhOCIsImlhdCI6MTU3Mzc1MzMwNH0.lbnYn9QUEI1HYqltuMYrGS2KT0swmtF-1X7QEdLMyHM"
	validKey := []byte("secret")

	h := JWTMiddleware(JWTConfig{
		AllowedRole: "admin",
		SigningKey:  validKey,
	})(handler)

	makeReq := func(token string) echo.Context {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()
		req.Header.Set(echo.HeaderAuthorization, "Bearer: "+token)
		c := e.NewContext(req, res)

		reqResult := h(c)
		if reqResult == nil {
			t.Errorf("Should have unathorized error, roles do not match")
		}
		return c
	}
	makeReq(secondToken)
	makeReq(initialToken)
}

func TestExpiredToken(t *testing.T) {
	e := echo.New()
	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}
	expiredToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImludmFsaWQtdG9rZW4iLCJ1c2VyX2lkIjoicmFuZG9tX2lkIiwicm9sZXMiOlsidXNlciJdLCJleHAiOjE1NzM3NTUzMDMsImp0aSI6ImVjYjM3OWUzLTk1NGUtNGE4NC1hNzNjLWE2OWNkMGZhOTIyZiIsImlhdCI6MTU3Mzc1NTMwMn0.PXh5_3Z5dpPI5zpqNIgXpUm4uc6HIuZMt-wddx44cDs"
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
		if reqResult == nil {
			t.Errorf("Should have unathorized error, token expired")
		}
		return c
	}
	makeReq(expiredToken)
}

func TestGetAuthTokenFromHeader(t *testing.T) {
	genContext := func(h string) echo.Context {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()
		if h != "" {
			req.Header.Set(echo.HeaderAuthorization, h)
		}
		c := e.NewContext(req, res)
		return c
	}

	_, err := getAuthTokenFromHeader(genContext("Beeaaaarer"))
	if err == nil {
		t.Errorf("Token is malformed, should have returned error")
	}

	_, err = getAuthTokenFromHeader(genContext("Bearer:"))
	if err != nil {
		t.Errorf("Token is just not supplied, should not error")
	}

	_, err = getAuthTokenFromHeader(genContext("not bearer"))
	if err == nil {
		t.Errorf("Token not set, should have returned error")
	}
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
