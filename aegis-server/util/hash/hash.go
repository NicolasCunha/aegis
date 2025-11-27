// Package hash provides secure password hashing utilities using HMAC-SHA256
// with salt and pepper for additional security.
package hash

import (
	"log"
	"os"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"crypto/rand"
)

var HASH_KEY = getHashKey()
const SALT_LENGTH = 16
const PEPPER_LENGTH = 16

// HashOutput represents the result of a password hashing operation,
// containing the hash value along with the salt and pepper used.
type HashOutput struct {
	Hash   string
	Salt   string
	Pepper string
}

// Hash generates a secure hash of the input string with randomly generated salt and pepper.
// The hash is computed using HMAC-SHA256 with a secret key.
//
// Parameters:
//   - input: The string to hash (typically a password)
//
// Returns:
//   - HashOutput containing the hash, salt, and pepper values
//
// Panics:
//   - If random number generation fails for salt or pepper
func Hash(input string) HashOutput {
	// Generate salt and pepper
	saltBytes := make([]byte, SALT_LENGTH)
	_, err := rand.Read(saltBytes)

	if err != nil {
		panic(err)
	}
	
	pepperBytes := make([]byte, PEPPER_LENGTH)
	_, err = rand.Read(pepperBytes)
	if err != nil {
		panic(err)
	}

	salt := hex.EncodeToString(saltBytes)
	pepper := hex.EncodeToString(pepperBytes)

	return HashWithSaltAndPepper(input, salt, pepper)
}

// HashWithSaltAndPepper generates a hash using the provided salt and pepper values.
// This is useful for validating passwords by recreating the hash with stored salt/pepper.
//
// Parameters:
//   - input: The string to hash
//   - salt: The salt value to use
//   - pepper: The pepper value to use
//
// Returns:
//   - HashOutput containing the computed hash along with the provided salt and pepper
func HashWithSaltAndPepper(input string, salt string, pepper string) HashOutput {
	// Combine input with salt, pepper, and secret key
	combined := input + salt + pepper

	// Hash SHA-256 signing it with the secret key
	hasher := hmac.New(sha256.New, []byte(HASH_KEY))
	hasher.Write([]byte(combined))

	// Get the sum
	hmacSum := hasher.Sum(nil)

	// Encode to hex string
	hash := hex.EncodeToString(hmacSum)

	// Return the Hash, salt, and pepper
	return HashOutput{
		Hash:   hash,
		Salt:   salt,
		Pepper: pepper,
	}
}

// Compare verifies if an input string matches a stored hash when using the same salt and pepper.
// This is used for password verification during authentication.
//
// Parameters:
//   - input: The string to verify (e.g., user-provided password)
//   - salt: The salt value from the stored hash
//   - pepper: The pepper value from the stored hash
//   - hash: The stored hash value to compare against
//
// Returns:
//   - true if the input generates the same hash, false otherwise
func Compare(input string, salt string, pepper string, hash string) bool {
	hashOutput := HashWithSaltAndPepper(input, salt, pepper)
	return hashOutput.Hash == hash
}

// getHashKey retrieves the HMAC secret key from the AEGIS_HASH_KEY environment variable.
// If not set, returns a default key with a warning. In production, always set this variable
// to a strong, random secret.
//
// Returns:
//   - The hash key string
func getHashKey() string {
	const HASH_KEY_ENV = "AEGIS_HASH_KEY"
	if key := os.Getenv(HASH_KEY_ENV); key != "" {
		return key
	}
	generatedHashKey := "DEFAULT_HASH_KEY"
	log.Printf("Warning: using default hash key '%s', consider setting the environment variable '%s'\n", generatedHashKey, HASH_KEY_ENV)
	return generatedHashKey
}