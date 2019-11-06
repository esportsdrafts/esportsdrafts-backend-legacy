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

// Structure comes from the offical JWT middleware in Echo
type (
	// JWTConfig defines the config for JWT middleware.
	JWTConfig struct {
		// SuccessHandler defines a function which is executed for a valid token.
		SuccessHandler JWTSuccessHandler

		// ErrorHandler defines a function which is executed for an invalid token.
		// It may be used to define a custom JWT error.
		ErrorHandler JWTErrorHandler

		// ErrorHandlerWithContext is almost identical to ErrorHandler, but it's passed the current context.
		ErrorHandlerWithContext JWTErrorHandlerWithContext

		// Signing key to validate token. Used as fallback if SigningKeys has length 0.
		// Required. This or SigningKeys.
		SigningKey interface{}
	}

	// JWTSuccessHandler defines a function which is executed for a valid token.
	JWTSuccessHandler func(echo.Context)

	// JWTErrorHandler defines a function which is executed for an invalid token.
	JWTErrorHandler func(error) error

	// JWTErrorHandlerWithContext is almost identical to JWTErrorHandler, but it's passed the current context.
	JWTErrorHandlerWithContext func(error, echo.Context) error
)

var (
	// ErrJWTMissing JWT Error
	ErrJWTMissing = echo.NewHTTPError(http.StatusBadRequest, "JWT token is missing or malformed")
)

// JWTMiddleware will check if token/cookie has correct signature,
// and if the allowed roles are in token
func JWTMiddleware(config JWTConfig) echo.MiddlewareFunc {
	if config.SigningKey == nil || config.SigningKey == "" {
		panic("JWT auth middleware requires signing secret")
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
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
	prefix := "Bearer"
	if strings.HasPrefix(headerContent, prefix) {
		runes := []rune(headerContent)
		if len(runes) <= len(prefix) {
			return "", fmt.Errorf("Auth header not found")
		}
		return strings.TrimSpace(string(runes[len(prefix):])), nil
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

// ReadAuthCookies get both header and signature from cookies
func readAuthCookies(ctx echo.Context) (string, error) {
	headerCookie, err := ctx.Cookie("header.payload")
	if err != nil {
		return "", err
	}

	signatureCookie, err := ctx.Cookie("signature")
	if err != nil {
		return "", err
	}
	return reconstructAuthToken(headerCookie.Value, signatureCookie.Value), nil
}

// ReconstructAuthToken join header and signature cookie values
func reconstructAuthToken(header, signature string) string {
	return header + "." + signature
}

func refreshToken() error {
	return nil
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
