package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/tarantool/go-tarantool/v2"
	"web.go/internal/home"
	"web.go/internal/notfound_navigate"
	"web.go/internal/rql"
	"web.go/internal/saveshare"
	"web.go/internal/tinyfy"
)

type paths_struct struct {
	HTML string // Exported field (starts with uppercase)
}

var paths = paths_struct{HTML: "/home/tinyurl/html/"}
var connsaveshare *tarantool.Connection

const prefix string = "https://example.com/"

// connectToTarantool initializes the connection to the Tarantool database.

type ChatPageData struct {
	InitialMessages template.HTML // Change type to template.HTML
}

// Main function
func main() {
	fmt.Println("prefix=", prefix)
	//SAVESHARE:
	http.HandleFunc("/saveshare", func(w http.ResponseWriter, r *http.Request) {
		saveshare.SavesharePageHandler(w, r, paths.HTML)
	})
	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		if rql.AllowAccess(w, r) { //check amount if requests
			saveshare.SubmitHandler(w, r, connsaveshare, prefix)
		} else {
			rql.AccessDenied(w, r, paths.HTML)
		}

	}) //handle form submissions
	http.HandleFunc("/link", func(w http.ResponseWriter, r *http.Request) {
		saveshare.LinkDisplayHandler(w, r, paths.HTML)
	}) //New handler for displaying the link
	http.HandleFunc("/ss/", func(w http.ResponseWriter, r *http.Request) {
		saveshare.SaveshareRedirectPageHandler(w, r, paths.HTML, connsaveshare)
	}) //Redirect handler

	//HOME:
	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		home.HomePageHandler(w, r, paths.HTML)
	}) // Home page handler

	//NOTFOUND_NAVIGATE:
	http.HandleFunc("/navigate", func(w http.ResponseWriter, r *http.Request) {
		notfound_navigate.NavigatePageHandler(w, r, paths.HTML)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		notfound_navigate.NotFoundHandler(w, r, paths.HTML)
	}) // This will handle all other routes

	//TINYFY:
	http.HandleFunc("/tinyfy", func(w http.ResponseWriter, r *http.Request) {
		tinyfy.TinyfyPageHandler(w, r, paths.HTML)
	})
	http.HandleFunc("/tinyfysubmit", func(w http.ResponseWriter, r *http.Request) {
		if rql.AllowAccess(w, r) { //check amount if requests
			tinyfy.TinyfySubmitHandler(w, r, paths.HTML)
		} else {
			rql.AccessDenied(w, r, paths.HTML)
		}
	})
	http.HandleFunc("/t/", func(w http.ResponseWriter, r *http.Request) {
		tinyfy.RedirectHandler(w, r)
	}) // Handle redirects at path "/t/"

	//CHAT:
	//http.HandleFunc("/get_new_admin_messages", func(w http.ResponseWriter, r *http.Request) {
	//	chat.GetNewAdminMessagesHandler(w, r, connwebchat)
	//})
	//http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
	//	chat.ChatHandler(w, r, connwebchat, paths.HTML)
	//})
	//http.HandleFunc("/send_message", func(w http.ResponseWriter, r *http.Request) {
	//	chat.SendMessageHandler(w, r, connwebchat)
	//})

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
