package tinyfy

import (
	"fmt"
	"strings"
)

func makeLinkFull(url string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}
	return url
}

// Get the original link by hash
func getLink(userid string, hash string) (string, string, error) {
	connheckfy, err := connectToTarantool()
	if err != nil {
		return "", "", fmt.Errorf("error connecting to Tarantool: %w", err)
	}
	defer connheckfy.Close()
	resp, err := connheckfy.Call("Get_link", []interface{}{userid, hash})
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
	link, ok := resp[0].(string)
	if !ok {
		return "", "", fmt.Errorf("invalid link type: %v", resp[0])
	}
	return link, statuscode, nil
}
