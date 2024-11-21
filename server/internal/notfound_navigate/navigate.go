package notfound_navigate

import (
	"fmt"
	"net/http"
	"path/filepath"
)

func NavigatePageHandler(w http.ResponseWriter, r *http.Request, path string) {
	//setUUIDCookie(w, r)                                              // Set cookie with UUID
	fmt.Println("Home page request: ", r.URL.Path)             // Debugging output
	http.ServeFile(w, r, filepath.Join(path, "navigate.html")) // home.html route
}
