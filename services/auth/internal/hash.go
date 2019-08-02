package internal

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	// ErrInvalidHash is returned if a hashed string contains wrong metadata
	ErrInvalidHash = errors.New("The encoded hash is not in the correct format")

	// ErrIncompatibleVersion is returned if the metadata indicates that the
	// argon hashing algorithm version is not the correct one.
	ErrIncompatibleVersion = errors.New("Incompatible version of argon2")
)

// Params holds hashing parameters for Argon2
type Params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

// SaltGenerator defines the functions to generate cryptographically secure
// random salts.
type SaltGenerator interface {
	Generate(n uint32) ([]byte, error)
}

// HashDecoder has required functions to take an encoded hash and decode it.
type HashDecoder interface {
	Decode(encodedHash string) (p *Params, salt, hash []byte, err error)
}

type RandReadGenerator struct{}
type Argon2HashDecoder struct{}

// GetDefaultHashingParams returns sane defaults for the hashing.
// Verified by unit tests. See 'internal/test_hash.go'.
func GetDefaultHashingParams() *Params {
	return &Params{
		memory:      64 * 1024,
		iterations:  3,
		parallelism: 2,
		saltLength:  16,
		keyLength:   32,
	}
}

// GenerateFromPassword generates a hashed string from a plain-text password
// using the specified parameters (usually 'GetDefaultHashingParams()').
func GenerateFromPassword(password string, p *Params) (hash string, err error) {
	// Generate a cryptographically secure random salt.
	salt, err := generateRandomBytes(p.saltLength)
	if err != nil {
		return "", err
	}

	hashed := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hashed)

	// Return a string using the standard encoded hash representation.
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, p.memory, p.iterations, p.parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// ComparePasswordAndHash takes a plain-text password and a hashed password,
// hashes the plain-text password using the parameters from the hashed one, and
// then compares them. Returns true if they are equal; otherwise false.
// -
// NOTE: Compares using a constant time comparison function. This means this
// 		 function is deterministically slow.
func ComparePasswordAndHash(password, encodedHash string) (match bool, err error) {
	// Extract the parameters, salt and derived key from the encoded password
	// hash.
	p, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Derive the key from the other password using the same parameters.
	otherHash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	return subtle.ConstantTimeCompare(hash, otherHash) == 1, nil
}

func decodeHash(encodedHash string) (p *Params, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	p = &Params{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.saltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.keyLength = uint32(len(hash))

	return p, salt, hash, nil
}
