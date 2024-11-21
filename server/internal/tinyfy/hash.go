package tinyfy

import (
	"fmt"

	"golang.org/x/crypto/argon2"
)

func getHash(userid string, link string) (string, string, error) {
	link = makeLinkFull(link)
	hashBytes, err := hashLink(link)
	if err != nil {
		return "", "", fmt.Errorf("error generating hash: %w", err)
	}
	hash := fmt.Sprintf("%x", hashBytes)
	connheckfy, err := connectToTarantool()
	if err != nil {
		return "", "", fmt.Errorf("error connecting to Tarantool: %w", err)
	}
	defer connheckfy.Close()
	resp, err := connheckfy.Call("Get_hash", []interface{}{userid, hash, link})
	if err != nil {
		return "", "", fmt.Errorf("error calling Lua function: %w", err)
	}
	if len(resp) < 2 {
		return "", "", fmt.Errorf("unexpected response from Lua function: %v", resp)
	}
	statuscode, ok := resp[1].(string)
	if !ok {
		return "", "", fmt.Errorf("invalid status code type: %v", resp[1])
	}
	return hash, statuscode, nil
}

func hashLink(link string) ([]byte, error) {
	iterations := uint32(1)
	memory := uint32(64 * 1024)
	parallelism := uint8(1)
	keyLength := uint32(8)
	salt := []byte("") // Keeping salt empty as per your requirement
	if link == "" {
		return nil, fmt.Errorf("link cannot be empty")
	}
	hash := argon2.IDKey([]byte(link), salt, iterations, memory, parallelism, keyLength)
	return hash, nil
}
