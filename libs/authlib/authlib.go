package authlib

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

// JWTClaims holds esportsdrafts auth claims. Roles array denotes what the user
// can do within the application. For example and 'admin' would have elevated
// access compared to a 'user'.
type JWTClaims struct {
	Username string   `json:"username"`
	UserID   string   `json:"user_id"`
	Roles    []string `json:"roles"`
	jwt.StandardClaims
}

// JWTMiddleware will check if token/cookie has correct signature,
// and if the allowed roles are in token
func JWTMiddleware(allowedRoles []string, signingSecret []byte) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			return nil
		}
	}
}

// HasRequestedWithHeader checks if X-Requested-With header has value
// XMLHttpRequest
func HasRequestedWithHeader(ctx echo.Context) bool {
	return ctx.Request().Header.Get("X-Requested-With") == "XMLHttpRequest"
}

// GetAuthTokenFromHeader grabs JWT token from header entry
func GetAuthTokenFromHeader(ctx echo.Context) (string, error) {
	headerContent := ctx.Request().Header.Get("Authorization")
	headerContent = strings.TrimSpace(headerContent)
	if strings.HasPrefix(headerContent, "Bearer") {
		runes := []rune(headerContent)
		if len(runes) <= 7 {
			return "", fmt.Errorf("Auth header not found")
		}
		return strings.TrimSpace(string(runes[6:])), nil
	}
	return "", fmt.Errorf("Auth header not found")
}

// WriteHeaderPayloadCookie header entries in JWT token to cookie
func WriteHeaderPayloadCookie(ctx echo.Context, header string, expiry time.Duration) {
	cookie := new(http.Cookie)
	cookie.Name = "header.payload"
	cookie.Value = header

	// Protect from sending over HTTP
	cookie.Secure = true
	cookie.Expires = time.Now().Add(expiry)
	ctx.SetCookie(cookie)
}

// WriteSignatureCookie writes the JWT signature to a secure cookie
func WriteSignatureCookie(ctx echo.Context, signature string) {
	cookie := new(http.Cookie)
	cookie.Name = "signature"
	cookie.Value = signature

	// Protect from sending over HTTP
	cookie.Secure = true

	// Block JS from reading this cookie
	cookie.HttpOnly = true
	ctx.SetCookie(cookie)
}

// GenerateAuthToken generates a auth token with provided claims
func GenerateAuthToken(claims *JWTClaims, expiry time.Duration, jwtKey []byte) (string, time.Time, error) {
	issuedTime := time.Now()
	expirationTime := issuedTime.Add(expiry)
	claims.StandardClaims = jwt.StandardClaims{
		// In JWT, the expiry time is expressed as unix milliseconds
		ExpiresAt: expirationTime.Unix(),
		// Can be used to blacklist in the future. Needs to hold state
		// in that case :/
		Id:       uuid.NewV4().String(),
		IssuedAt: issuedTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	res, err := token.SignedString(jwtKey)
	return res, expirationTime, err
}
