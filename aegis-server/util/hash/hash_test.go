package hash

import (
	"os"
	"testing"
)

// TestHash tests the Hash function with random salt and pepper
func TestHash(t *testing.T) {
	input := "password123"
	
	result := Hash(input)
	
	// Verify all fields are populated
	if result.Hash == "" {
		t.Error("Hash should not be empty")
	}
	if result.Salt == "" {
		t.Error("Salt should not be empty")
	}
	if result.Pepper == "" {
		t.Error("Pepper should not be empty")
	}
	
	// Verify salt and pepper have expected length (32 hex chars = 16 bytes)
	if len(result.Salt) != SALT_LENGTH*2 {
		t.Errorf("Salt length should be %d, got %d", SALT_LENGTH*2, len(result.Salt))
	}
	if len(result.Pepper) != PEPPER_LENGTH*2 {
		t.Errorf("Pepper length should be %d, got %d", PEPPER_LENGTH*2, len(result.Pepper))
	}
	
	// Verify hash is deterministic - same input with same salt/pepper produces same hash
	result2 := HashWithSaltAndPepper(input, result.Salt, result.Pepper)
	if result.Hash != result2.Hash {
		t.Error("Hash should be deterministic with same salt and pepper")
	}
}

// TestHashUniqueness tests that Hash generates unique salt and pepper each time
func TestHashUniqueness(t *testing.T) {
	input := "password123"
	
	result1 := Hash(input)
	result2 := Hash(input)
	
	// Same input should produce different hashes due to random salt/pepper
	if result1.Hash == result2.Hash {
		t.Error("Hash should be unique due to random salt and pepper")
	}
	if result1.Salt == result2.Salt {
		t.Error("Salt should be unique")
	}
	if result1.Pepper == result2.Pepper {
		t.Error("Pepper should be unique")
	}
}

// TestHashWithSaltAndPepper tests hashing with provided salt and pepper
func TestHashWithSaltAndPepper(t *testing.T) {
	input := "password123"
	salt := "a1b2c3d4e5f6"
	pepper := "1a2b3c4d5e6f"
	
	result := HashWithSaltAndPepper(input, salt, pepper)
	
	if result.Hash == "" {
		t.Error("Hash should not be empty")
	}
	if result.Salt != salt {
		t.Errorf("Expected salt %s, got %s", salt, result.Salt)
	}
	if result.Pepper != pepper {
		t.Errorf("Expected pepper %s, got %s", pepper, result.Pepper)
	}
}

// TestHashWithSaltAndPepper_Deterministic tests that hash is deterministic
func TestHashWithSaltAndPepper_Deterministic(t *testing.T) {
	input := "password123"
	salt := "a1b2c3d4e5f6"
	pepper := "1a2b3c4d5e6f"
	
	result1 := HashWithSaltAndPepper(input, salt, pepper)
	result2 := HashWithSaltAndPepper(input, salt, pepper)
	
	if result1.Hash != result2.Hash {
		t.Error("Hash should be deterministic with same input, salt, and pepper")
	}
}

// TestHashWithSaltAndPepper_DifferentInputs tests different inputs produce different hashes
func TestHashWithSaltAndPepper_DifferentInputs(t *testing.T) {
	salt := "a1b2c3d4e5f6"
	pepper := "1a2b3c4d5e6f"
	
	result1 := HashWithSaltAndPepper("password1", salt, pepper)
	result2 := HashWithSaltAndPepper("password2", salt, pepper)
	
	if result1.Hash == result2.Hash {
		t.Error("Different inputs should produce different hashes")
	}
}

// TestHashWithSaltAndPepper_DifferentSalt tests different salt produces different hashes
func TestHashWithSaltAndPepper_DifferentSalt(t *testing.T) {
	input := "password123"
	pepper := "1a2b3c4d5e6f"
	
	result1 := HashWithSaltAndPepper(input, "salt1", pepper)
	result2 := HashWithSaltAndPepper(input, "salt2", pepper)
	
	if result1.Hash == result2.Hash {
		t.Error("Different salts should produce different hashes")
	}
}

// TestHashWithSaltAndPepper_DifferentPepper tests different pepper produces different hashes
func TestHashWithSaltAndPepper_DifferentPepper(t *testing.T) {
	input := "password123"
	salt := "a1b2c3d4e5f6"
	
	result1 := HashWithSaltAndPepper(input, salt, "pepper1")
	result2 := HashWithSaltAndPepper(input, salt, "pepper2")
	
	if result1.Hash == result2.Hash {
		t.Error("Different peppers should produce different hashes")
	}
}

// TestCompare_ValidPassword tests successful password verification
func TestCompare_ValidPassword(t *testing.T) {
	input := "password123"
	hashOutput := Hash(input)
	
	result := Compare(input, hashOutput.Salt, hashOutput.Pepper, hashOutput.Hash)
	
	if !result {
		t.Error("Compare should return true for valid password")
	}
}

// TestCompare_InvalidPassword tests failed password verification
func TestCompare_InvalidPassword(t *testing.T) {
	input := "password123"
	wrongInput := "wrongpassword"
	hashOutput := Hash(input)
	
	result := Compare(wrongInput, hashOutput.Salt, hashOutput.Pepper, hashOutput.Hash)
	
	if result {
		t.Error("Compare should return false for invalid password")
	}
}

// TestCompare_WrongSalt tests password verification fails with wrong salt
func TestCompare_WrongSalt(t *testing.T) {
	input := "password123"
	hashOutput := Hash(input)
	
	result := Compare(input, "wrongsalt", hashOutput.Pepper, hashOutput.Hash)
	
	if result {
		t.Error("Compare should return false with wrong salt")
	}
}

// TestCompare_WrongPepper tests password verification fails with wrong pepper
func TestCompare_WrongPepper(t *testing.T) {
	input := "password123"
	hashOutput := Hash(input)
	
	result := Compare(input, hashOutput.Salt, "wrongpepper", hashOutput.Hash)
	
	if result {
		t.Error("Compare should return false with wrong pepper")
	}
}

// TestCompare_EmptyPassword tests comparison with empty password
func TestCompare_EmptyPassword(t *testing.T) {
	hashOutput := Hash("")
	
	result := Compare("", hashOutput.Salt, hashOutput.Pepper, hashOutput.Hash)
	
	if !result {
		t.Error("Compare should handle empty passwords correctly")
	}
}

// TestGetHashKey_DefaultValue tests default hash key is used when env var not set
func TestGetHashKey_DefaultValue(t *testing.T) {
	// Unset the environment variable
	os.Unsetenv("AEGIS_HASH_KEY")
	
	key := getHashKey()
	
	if key != "DEFAULT_HASH_KEY" {
		t.Errorf("Expected default hash key 'DEFAULT_HASH_KEY', got '%s'", key)
	}
}

// TestGetHashKey_CustomValue tests custom hash key from environment variable
func TestGetHashKey_CustomValue(t *testing.T) {
	customKey := "my_custom_secret_key"
	os.Setenv("AEGIS_HASH_KEY", customKey)
	defer os.Unsetenv("AEGIS_HASH_KEY")
	
	key := getHashKey()
	
	if key != customKey {
		t.Errorf("Expected custom hash key '%s', got '%s'", customKey, key)
	}
}

// TestHashWithCustomHashKey tests that different hash keys produce different hashes
func TestHashWithCustomHashKey(t *testing.T) {
	input := "password123"
	salt := "a1b2c3d4e5f6"
	pepper := "1a2b3c4d5e6f"
	
	// Use default key
	os.Unsetenv("AEGIS_HASH_KEY")
	HASH_KEY = getHashKey()
	result1 := HashWithSaltAndPepper(input, salt, pepper)
	
	// Use custom key
	os.Setenv("AEGIS_HASH_KEY", "different_key")
	HASH_KEY = getHashKey()
	result2 := HashWithSaltAndPepper(input, salt, pepper)
	defer os.Unsetenv("AEGIS_HASH_KEY")
	
	if result1.Hash == result2.Hash {
		t.Error("Different hash keys should produce different hashes")
	}
	
	// Reset to default for other tests
	os.Unsetenv("AEGIS_HASH_KEY")
	HASH_KEY = getHashKey()
}
