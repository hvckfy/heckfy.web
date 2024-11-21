package saveshare

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

func if_valid_json(json_to_verify json_data, hash string) bool { //this function is made to provied that data inside of json is valid or not
	return json_to_verify.Hash == hash
}

// Function to unpad the input data
func unpad(data []byte) ([]byte, error) { //This function unpad removes PKCS#7 padding from the input byte slice data, returning the unpadded data and an error if the padding size is invalid.
	padding := data[len(data)-1]
	if int(padding) > aes.BlockSize {
		return nil, fmt.Errorf("padding size is invalid")
	}
	return data[:len(data)-int(padding)], nil
}

// Function to decrypt the encrypted data
func decrypt(encrypted_data string, key []byte) (json_data, error) { // Define the decrypt function with parameters for encrypted data and key
	ciphertext, _ := base64.StdEncoding.DecodeString(encrypted_data) // Decode the base64-encoded string into a byte slice

	block, err := aes.NewCipher(key) // Create a new AES cipher block using the provided key
	if err != nil {                  // Check for errors in cipher creation
		return json_data{}, err // Return an empty json_data struct and the error if cipher creation fails
	}

	if len(ciphertext) < aes.BlockSize { // Check if the length of the ciphertext is less than the block size
		return json_data{}, fmt.Errorf("ciphertext too short") // Return an error indicating the ciphertext is too short
	}

	iv := ciphertext[:aes.BlockSize]        // Extract the initialization vector (IV) from the first block of the ciphertext
	ciphertext = ciphertext[aes.BlockSize:] // Remove the IV from the ciphertext, leaving only the encrypted data

	mode := cipher.NewCBCDecrypter(block, iv) // Create a new CBC (Cipher Block Chaining) decryption mode with the block and IV
	mode.CryptBlocks(ciphertext, ciphertext)  // Decrypt the ciphertext in place, modifying the ciphertext slice to hold the decrypted data

	unpaddedData, err := unpad(ciphertext) // Unpad the decrypted data to remove any padding added during encryption
	if err != nil {                        // Check for errors during unpadding
		return json_data{}, err // Return an empty json_data struct and the error if unpadding fails
	}

	var decryptedData json_data                        // Declare a variable to hold the unmarshalled JSON data
	err = json.Unmarshal(unpaddedData, &decryptedData) // Unmarshal the unpadded data into the decryptedData struct
	if err != nil {                                    // Check for errors during unmarshalling
		return json_data{}, err // Return an empty json_data struct and the error if unmarshalling fails
	}

	return decryptedData, nil // Return the successfully decrypted JSON data and a nil error
}
