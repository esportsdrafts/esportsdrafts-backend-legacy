package internal

import (
	"encoding/base64"
	"fmt"
	"github.com/esportsdrafts/esportsdrafts/libs/log"
	"math"
	"testing"
	"time"

	"golang.org/x/crypto/argon2"
)

func TestDefaultParameters(t *testing.T) {
	defaults := GetDefaultHashingParams()

	if defaults.memory < (32 * 1024) {
		t.Errorf("Too low hashing memory")
	}

	if defaults.iterations < 3 {
		t.Errorf("To few number of hashing iterations")
	}

	if defaults.parallelism < 2 {
		t.Errorf("Hashing defaults using too few threads")
	}

	if defaults.saltLength < 16 {
		t.Errorf("Insecure salt length")
	}

	if defaults.keyLength < 32 {
		t.Errorf("Insecure key length")
	}
}

func TestGenNullPassword(t *testing.T) {
	params := GetDefaultHashingParams()
	hashedPassword, _ := GenerateFromPassword("Ewdf3UKB8BGB1gWjkwvCWf6FZ3ZcYi8YqHVEDTRFZaYjqNXeGrHQH476kuEs2FMdmgPmY9RNjDjfACeuh1pcIA66GGCZ8Xu0hGBldr3s87Yc4iuwJCncEVJy", params)
	log.GetLogger().Infof("%s", hashedPassword)
}

// Argon2 hashing has to take a sufficient amount of time to not be
// easily cracked. However, spending too much time would open us up to
// some nasty DDos attacks, and generally slow API.
// What this test tries to do is show the different runtimes for a set
// of hashing parameters and detect if its too low/high.
// NOTE: This is super platform dependent so test might fail randomly
func TestGenerateHashingTimeDefaults(t *testing.T) {
	params := GetDefaultHashingParams()
	defer timeTrack(t, time.Now(), "Hashing password with default parameters", 50, 200)
	GenerateFromPassword("random_password_123213211", params)
}

func TestCompareTimeCorrectDefaults(t *testing.T) {
	params := GetDefaultHashingParams()
	clearText := "random_password_123213211"
	hashed, _ := GenerateFromPassword(clearText, params)
	defer timeTrack(t, time.Now(), "Compare correct with default parameters", 50, 200)
	ComparePasswordAndHash(clearText, hashed)
}

func TestCompareTimeNotMatchingDefaults(t *testing.T) {
	params := GetDefaultHashingParams()
	clearText := "random_password_123213211"
	hashed, _ := GenerateFromPassword(clearText, params)
	defer timeTrack(t, time.Now(), "Compare not matching with default parameters", 50, 200)
	ComparePasswordAndHash("SomE_other_r4nd0m_password", hashed)
}

func TestCompareTimeConstantDefaults(t *testing.T) {
	params := GetDefaultHashingParams()
	testStrings := []string{
		"random_password_123213211",
		"SomE_other_r4nd0m_password",
		"as",
		"-0-------",
		"SomE_other_r4nd0m_password",
		"random_password_123213211",
	}
	hashed, _ := GenerateFromPassword(testStrings[0], params)

	timings := []float64{}
	for _, password := range testStrings {
		start := time.Now()
		ComparePasswordAndHash(password, hashed)
		elapsed := time.Since(start)
		timings = append(timings, elapsed.Minutes()*1000)
	}

	mean := mean(timings)
	maxVariance := 10.0
	for _, timing := range timings {
		if timing > (mean+maxVariance) || timing < (mean-maxVariance) {
			t.Errorf("ComparePasswordAndHash has inconsistent compare time")
		}
	}
}

func TestDecodeHashErrors(t *testing.T) {
	p := GetDefaultHashingParams()
	extraEntries := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s$pp=23$extra=values",
		argon2.Version, p.memory, p.iterations, p.parallelism, "bogus", "bogus")

	pRes, s, h, err := decodeHash(extraEntries)
	if err != ErrInvalidHash {
		t.Errorf("Invalid number of entries did not throw correct error")
	}
	if pRes != nil || s != nil || h != nil {
		t.Errorf("DecodeHash returned bogus data with incorrect number of entries")
	}

	insufficientEntries := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s",
		argon2.Version, p.memory, p.iterations, p.parallelism, "bogus")

	pRes, s, h, err = decodeHash(insufficientEntries)
	if err != ErrInvalidHash {
		t.Errorf("Invalid number of entries did not throw correct error")
	}
	if pRes != nil || s != nil || h != nil {
		t.Errorf("DecodeHash returned bogus data with incorrect number of entries")
	}

	invalidVersionPosition := fmt.Sprintf("$argon2id$m=%d,t=%d,$v=%d,p=%d$%s$%s",
		argon2.Version, p.memory, p.iterations, p.parallelism, "bogus", "bogus")

	pRes, s, h, err = decodeHash(invalidVersionPosition)
	if err == nil {
		t.Errorf("Decode with invalid version position did not error out")
	}
	if err == ErrInvalidHash || err == ErrIncompatibleVersion {
		t.Errorf("Decode with invalid version position returned wrong error")
	}
	if pRes != nil || s != nil || h != nil {
		t.Errorf("DecodeHash returned bogus results with invalid Argon2 version")
	}

	invalidVersion := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		13, p.memory, p.iterations, p.parallelism, "bogus", "bogus")

	pRes, s, h, err = decodeHash(invalidVersion)
	if err == nil {
		t.Errorf("Decode with invalid version did not return an error")
	}
	if err != ErrIncompatibleVersion {
		t.Errorf("Decode with invalid version returned wrong error")
	}
	if pRes != nil || s != nil || h != nil {
		t.Errorf("DecodeHash returned bogus results with invalid Argon2 version")
	}

	invalidArgumentPositions := fmt.Sprintf("$argon2id$v=%d$t=%d,p=%d,m=%d$%s$%s",
		argon2.Version, p.memory, p.iterations, p.parallelism, "bogus", "bogus")

	// TODO: Simplify the code with a map of arguments and loop over it
	pRes, s, h, err = decodeHash(invalidArgumentPositions)
	if err == nil {
		t.Errorf("Decode with invalid positions did not return an error")
	}
	if pRes != nil || s != nil || h != nil {
		t.Errorf("DecodeHash returned bogus results with argument positions")
	}

	invalidB64Salt := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, p.memory, p.iterations, p.parallelism, "Ym9ndXNfc2FsdA==", "bogus")

	pRes, s, h, err = decodeHash(invalidB64Salt)
	if err == nil {
		t.Errorf("Decode with invalid positions did not return an error")
	}
	if pRes != nil || s != nil || h != nil {
		t.Errorf("DecodeHash returned bogus results with argument positions")
	}

	salt, err := generateRandomBytes(p.saltLength)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	invalidB64Hash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, p.memory, p.iterations, p.parallelism, b64Salt, "bogus")

	pRes, s, h, err = decodeHash(invalidB64Hash)
	if err == nil {
		t.Errorf("Decode invalid salt b64 string did not return error")
	}
	if pRes != nil || s != nil || h != nil {
		t.Errorf("DecodeHash returned bogus results with invalid salt b64")
	}
}

func TestGenFromPasswordError(t *testing.T) {
	hash, err := fGenFromPassword("bogus_password", GetDefaultHashingParams(), &ErrorMockGenerator{})
	if hash != "" {
		t.Errorf("GenFromPassword returned value on salt gen failure")
	}
	if err == nil {
		t.Errorf("GenFromPassword did not return error on salt gen failure")
	}
}

func TestComparePasswordAndHashError(t *testing.T) {
	encoded, _ := GenerateFromPassword("bogus_password", GetDefaultHashingParams())
	match, err := fComparePasswordAndHash("bogus_password", encoded, &ErrorMockDecoder{})
	if match {
		t.Errorf("ComparePasswordAndHash match on decode failure")
	}
	if err == nil {
		t.Errorf("ComparePasswordAndHash did not return error on salt gen failure")
	}
}

//
// Testing Helpers
//
func mean(n []float64) float64 {
	total := 0.0
	for _, v := range n {
		total += v
	}
	return math.Round(total / float64(len(n)))
}

// Track timing of function calls.
//
// Example Usage:
// ```
// defer timeTrack(t, time.Now(), "Arbitrary name", 10, 200)
// ```
// where `t` is `testing.T`
//
// Bounds are measured in milliseconds
func timeTrack(t *testing.T, start time.Time, name string, lowerBound float64, upperBound float64) {
	elapsed := time.Since(start)
	if (elapsed.Seconds() * 1000) > upperBound {
		t.Errorf("'%s' higher than upper limit (limit: %f, got: %s)", name, upperBound, elapsed)
	}
	if (elapsed.Seconds() * 1000) < lowerBound {
		t.Errorf("'%s' lower than lower limit (limit: %f, got: %s)", name, upperBound, elapsed)
	}
}
