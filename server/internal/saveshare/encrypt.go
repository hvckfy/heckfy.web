package saveshare

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"
)

// Function to create a SHA-256 hash of the input data
func create_hash(data string) string { //create a hash that we gonna use later
	hash := sha256.New()                     //lets say user paste some "data" and want to get link
	hash.Write([]byte(data))                 //he gets link a type of: https://domain.ru/"hash"&key="key"
	return hex.EncodeToString(hash.Sum(nil)) //key explain u can find later->
}

// Function to create a random AES key
func create_random_key() ([]byte, error) { //this function gonna create random key that actually a 16-byte data
	key := make([]byte, 16)                                  // AES requires a 16-byte key	//creates a byte slice (an array of bytes) with a length of 16, initialized with zero values.
	if _, err := io.ReadFull(rand.Reader, key); err != nil { //This code attempts to fill the key byte slice with random data from rand.Reader,
		return nil, err //and if it fails, it returns an error;
	}
	return key, nil //otherwise in returns key with err=<nil>
}

// Function to pad the input data
func pad(data []byte) []byte { //This function pad adds PKCS#7 padding to the input byte slice data to ensure its length is a multiple of the AES block size, allowing for proper encryption.
	padding := aes.BlockSize - len(data)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// Function to encrypt the JSON data
func encrypt(not_encrypted_json_data json_data, key []byte) (string, error) { // Function to encrypt JSON data using AES
	block, err := aes.NewCipher(key) // Create a new AES cipher block using the provided key
	if err != nil {                  // Check if there was an error creating the cipher
		return "", err // Return an empty string and the error if creation failed
	}

	jsonBytes, err := json.Marshal(not_encrypted_json_data) // Convert the JSON data to a byte slice
	if err != nil {                                         // Check if there was an error during marshaling
		return "", err // Return an empty string and the error if marshaling failed
	}

	paddedText := pad(jsonBytes) // Apply PKCS#7 padding to the JSON byte slice

	ciphertext := make([]byte, aes.BlockSize+len(paddedText)) // Allocate space for the ciphertext (IV + padded text)
	iv := ciphertext[:aes.BlockSize]                          // Extract the IV (initialization vector) from the beginning of the ciphertext
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {   // Fill the IV with random bytes
		return "", err // Return an empty string and the error if reading random bytes failed
	}

	mode := cipher.NewCBCEncrypter(block, iv)                // Create a new CBC encrypter with the cipher block and IV
	mode.CryptBlocks(ciphertext[aes.BlockSize:], paddedText) // Encrypt the padded text and store it in the ciphertext

	return base64.StdEncoding.EncodeToString(ciphertext), nil // Encode the ciphertext to a base64 string and return it
}
