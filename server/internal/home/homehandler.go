package home

import (
	"fmt"
	"net/http"
	"path/filepath"
)

// Home page handler
func HomePageHandler(w http.ResponseWriter, r *http.Request, path string) {
	//setUUIDCookie(w, r)                                          // Set cookie with UUID
	fmt.Println("Home page request: ", r.URL.Path)         // Debugging output
	http.ServeFile(w, r, filepath.Join(path, "home.html")) // home.html route
}
