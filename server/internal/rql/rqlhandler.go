package rql

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/tarantool/go-tarantool/v2"
)

func AllowAccess(w http.ResponseWriter, r *http.Request) bool {

	var connrql *tarantool.Connection

	var err error

	userIP := r.RemoteAddr
	if forwardedIP := r.Header.Get("X-Forwarded-For"); forwardedIP != "" {
		userIP = forwardedIP
	}

	fmt.Printf("Requests limiter has been called for IP: %s\n", userIP)

	connrql, err = RQLConnectToTarantool()
	if err != nil {
		fmt.Println("Failed to connect to Tarantool:", err)
		return false
	}
	defer connrql.Close() // Ensure the connection is closed when done

	fmt.Println("Successfully connected to Tarantool")

	// Call the Tarantool function and handle the result
	status, err := connrql.Call("access", []interface{}{userIP})
	if err != nil {
		fmt.Println("Failed to call Tarantool function:", err)
		return false
	}

	// Unpack the response
	if len(status) < 2 {
		fmt.Println("Unexpected response format:", status)
		return false
	}

	statusCode, ifAccess := status[0].(string), status[1].(bool)

	// Log the results
	fmt.Printf("Response from Tarantool: status code %s, access granted: %v\n", statusCode, ifAccess)
	// Implement your logic based on the response
	return ifAccess
}

func AccessDenied(w http.ResponseWriter, r *http.Request, path string) {
	fmt.Println("Home page request: ", r.URL.Path)                 // Debugging output
	http.ServeFile(w, r, filepath.Join(path, "accessdenied.html")) // home.html route
}
