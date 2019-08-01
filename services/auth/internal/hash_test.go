package internal

import (
	"testing"
	"time"
	"log"
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

// Argon2 hashing has to take a sufficient amount of time to not be
// easily cracked. However, spending too much time would open us up to
// some nasty DDos attacks, and generally slow API.
// What this test tries to do is show the different runtimes for a set
// of hashing parameters and detect if its too low/high.
// NOTE: This is super platform dependent so test might fail randomly
func TestGenerateHashingTimeDefaults(t *testing.T) {
	params := GetDefaultHashingParams()
	defer timeTrack(t, time.Now(), "Hashing passwird with default parameters", 50, 200)
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

// Bounds are measured in milliseconds
func timeTrack(t *testing.T, start time.Time, name string, lowerBound float64, upperBound float64) {
	elapsed := time.Since(start)
	if (elapsed.Seconds() * 1000) > upperBound {
		t.Errorf("'%s' higher than upper limit (limit: %f, got: %s)", name, upperBound, elapsed)
	}
	if (elapsed.Seconds() * 1000) < lowerBound {
		t.Errorf("'%s' lower than lower limit (limit: %f, got: %s)", name, upperBound, elapsed)
	}
	log.Printf("%s took %s", name, elapsed)
}
