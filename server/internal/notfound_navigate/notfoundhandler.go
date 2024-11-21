package notfound_navigate

import (
	"fmt"
	"net/http"
	"path/filepath"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request, path string) {
	// Create a set with file names
	fileSet := map[string]struct{}{
		"chat.html":   {},
		"chat.js":     {},
		"index.html":  {},
		"openchat.js": {},
	}
	// Check if the requested file is in the set
	if _, exists := fileSet[filepath.Base(r.URL.Path)]; exists {
		// If the file exists in the set, redirect to its handler
		http.ServeFile(w, r, filepath.Join(path, filepath.Base(r.URL.Path)))
	} else {
		fmt.Println("Redirect from notfound req: ", r.URL.Path)    // Debugging output
		http.ServeFile(w, r, filepath.Join(path, "navigate.html")) // home.html route
	}
}
