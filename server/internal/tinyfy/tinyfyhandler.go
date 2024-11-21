package tinyfy

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
)

func serveResponse(w http.ResponseWriter, shortenedLink string, path string) {
	htmlFilePath := filepath.Join(path, "response.html") // Path to the response HTML file
	data := struct {
		ShortenedLink string
	}{
		ShortenedLink: shortenedLink,
	}
	// Parsing the response HTML file
	t, err := template.ParseFiles(htmlFilePath)
	if err != nil {
		fmt.Println("Error reading response template:", err) // Output error to console
		http.Error(w, "Error reading response template", http.StatusInternalServerError)
		return
	}
	// Debugging output to check the value of ShortenedLink
	fmt.Println("Shortened Link:", shortenedLink)
	w.Header().Set("Content-Type", "text/html")
	if err := t.Execute(w, data); err != nil {
		fmt.Println("Error executing template:", err) // Output error to console
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}

func TinyfyPageHandler(w http.ResponseWriter, r *http.Request, path string) {
	//setUUIDCookie(w, r)                                      // Set cookie with UUID
	htmlFilePath := filepath.Join(path, "tinyfy.html") // Path to the HTML file
	http.ServeFile(w, r, htmlFilePath)                 // Serve the file
}

func TinyfySubmitHandler(w http.ResponseWriter, r *http.Request, path string) {
	userid := "0"               // test userid=0
	link := r.FormValue("link") // Get the link from the POST request
	if len(link) >= 10*1024 {
		http.Error(w, "Too many symbols, maximum allowed size is 10MB", http.StatusBadRequest)
		return // end
	}
	hash, statuscode, err := getHash(userid, link)
	if err != nil {
		fmt.Println("Statuscode:", statuscode)
		fmt.Println("Error generating hash:", err)
		http.Error(w, "Error generating hash", http.StatusInternalServerError)
		return
	}
	// Prepare the shortened link
	postfix := "t/" //replace with redirect domain
	shortenedLink := postfix + hash
	// Send the response with the shortened link
	serveResponse(w, shortenedLink, path)
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	//setUUIDCookie(w, r)                           // Set cookie with UUID
	hash := strings.TrimPrefix(r.URL.Path, "/t/") // Remove the leading slash and "/t/"
	userid := "0"                                 // Here you can replace with the actual userid if needed
	fmt.Println("Redirecting for hash:", hash)    // Debugging output

	link, _, err := getLink(userid, hash)
	if err != nil {
		fmt.Println("Error retrieving link:", err) // Debugging output
		http.Error(w, "Link not found", http.StatusNotFound)
		return
	}

	// Check if link is equal to hash
	if link == hash {
		fmt.Println("Link is equal to hash, redirecting to 404 page") // Debugging output
		http.Redirect(w, r, "http://heckfy.ru/404", http.StatusFound)
		return
	}

	fmt.Println("Redirecting to link:", link) // Debugging output
	http.Redirect(w, r, link, http.StatusFound)
}
