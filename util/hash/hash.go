package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"crypto/rand"
)

const Hash_KEY = "some_secret_key" // TODO: move to config/env variables
const SALT_LENGTH = 16
const PEPPER_LENGTH = 16

/* Hash output structure */
type HashOutput struct {
	Hash   string
	Salt   string
	Pepper string
}

/* Hash with generated salt and pepper */
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

/* Hash with provided salt and pepper */
func HashWithSaltAndPepper(input string, salt string, pepper string) HashOutput {
	// Combine input with salt, pepper, and secret key
	combined := input + salt + pepper

	// Hash SHA-256 signing it with the secret key
	hasher := hmac.New(sha256.New, []byte(Hash_KEY))
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

/* Compare two Hashes */
func Compare(input string, salt string, pepper string, Hash string) bool {
	HashOutput := HashWithSaltAndPepper(input, salt, pepper)
	return HashOutput.Hash == Hash
}