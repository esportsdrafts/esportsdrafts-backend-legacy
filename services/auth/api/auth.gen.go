// Package auth provides primitives to interact the openapi HTTP API.
//
// This is an autogenerated file, any edits which you make here will be lost!
package auth

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"strings"
)

// Account defines component schema for Account.
type Account struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

// AuthClaim defines component schema for AuthClaim.
type AuthClaim struct {
	Claim    string  `json:"claim"`
	Code     *string `json:"code,omitempty"`
	Password *string `json:"password,omitempty"`
	Username *string `json:"username,omitempty"`
}

// Error defines component schema for Error.
type Error struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

// JWT defines component schema for JWT.
type JWT struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int32  `json:"expires_in"`
	IdToken     string `json:"id_token"`
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Authenticate a user returning a JWT for future operations and set session token for browsers (POST /v1/auth/auth)
	PerformAuth(ctx echo.Context) error
	// Create a new account (POST /v1/auth/register)
	CreateAccount(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// PerformAuth converts echo context to params.
func (w *ServerInterfaceWrapper) PerformAuth(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PerformAuth(ctx)
	return err
}

// CreateAccount converts echo context to params.
func (w *ServerInterfaceWrapper) CreateAccount(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.CreateAccount(ctx)
	return err
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router runtime.EchoRouter, si ServerInterface) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST("/v1/auth/auth", wrapper.PerformAuth)
	router.POST("/v1/auth/register", wrapper.CreateAccount)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/7xVTY/bRgz9KwO2twiWvc5Jp26DtmhOBZoih4URTCRKnkaamZKc9RoL//eCI3/IsRbY",
	"Amku1mA++B75HulnqMMQg0cvDNUzcL3FweblfV2H5EWXkUJEEof5AAfrel3IPiJUwELOd3AoIFrmXaBm",
	"9jAxkrcDzhweCiD8JznCBqqHI8DkxSTy5lDAfZLtu9664ZZbfdpGnwaNdYrx5hyhgLvWgiJ63MGmuGVa",
	"hwb/z/xGkprJL0SBZrI4EmgDDVagAudlfQdnqs4LdkiKOiCz7V4DqjEv9xX9/ccPt9i2rpH5k4Qv6Gfz",
	"xKfoCPmTmx5PGLnmxcdfUbqCmjy8wtgcDlkSz2lQhg9gY+xdbcUFX/7NwauGzrchAzrpFRF/tV4s741N",
	"soUCHpHYBQ8VrBbLxVJ5hojeRgcVrPOWyivbXIPycVXqw/yTSxQ4N4IWKgP/3kAFEUkVuh8hNDNk+Tk0",
	"+1FDLzi2zw3hc6Pp6kfCFir4obx0Ynlsw/Li9FyFBrkmF2XMRA/RyzGyya5iY30zLk1jxcK04kIJswQc",
	"g+dR7rvl8puxVT/N8Hz/8UOWYcJ11DnfbG3q5ZtRGBtqhkTy+BSxFmwMHu8UwGkYLO2va4nGGu1lQyiJ",
	"vPOdsUZzaAOZNkkiNGcfjPVmFMPIfE4t3/1MYcdIrH1ru9G6apWNQp8tRtg5FqSXbVYTWsHTOP6vRsMn",
	"O8QeJ4N7/P60DSw6rxZ1GKYTttJm2f8ptPTdH2+Zd8s8NC/jDbq3w5p+W1G/WquGr7Tykf+ckccj06BY",
	"1/MrPLvSz3wQTnmqtKnv92YsXfNdnfaX/+LDzs/a7F3mY6zxuDP2rOjX9tC/GgpNql8YeMXVVqTwucfh",
	"zWkWMpKOO6geniGRCr4ViVyVpY1uge04GRcNPsJhc/g3AAD//+prFHMCCAAA",
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file.
func GetSwagger() (*openapi3.Swagger, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromData(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("error loading Swagger: %s", err)
	}
	return swagger, nil
}

