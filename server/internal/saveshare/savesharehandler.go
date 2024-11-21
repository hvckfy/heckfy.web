package saveshare

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"strconv"

	"github.com/tarantool/go-tarantool/v2"
)

func put_data_in_json(hash string, not_encrypted_data string) json_data { //this function is not neccessory
	return json_data{ //but its good idea to use functions as this
		Hash: hash,               //later i would not have to always declare everything is structure
		Data: not_encrypted_data, //use this instead of u see next lines:
	}
}

func create_link(urlprefix string, hash string, key []uint8) string { //this function are made to create a data link for a person after he send request
	link := urlprefix + "ss/" + hash + "&key=" + hex.EncodeToString(key)
	fmt.Println(link)
	return link
}

func SavesharePageHandler(w http.ResponseWriter, r *http.Request, path string) {
	//cookiework.setUUIDCookie(w, r) // Set cookie with UUID
	http.ServeFile(w, r, filepath.Join(path, "saveshare.html"))
}

func LinkDisplayHandler(w http.ResponseWriter, r *http.Request, path string) {
	encodedLink := r.URL.Query().Get("link") // Get the encoded link from the query parameters
	if encodedLink == "" {
		http.Error(w, "No link provided", http.StatusBadRequest)
		return
	}

	link, err := url.QueryUnescape(encodedLink) // Decode the URL
	if err != nil {
		http.Error(w, "Error decoding link", http.StatusBadRequest)
		return
	}

	// Read the HTML template
	htmlPath := filepath.Join(path, "link_display.html")
	htmlContent, err := ioutil.ReadFile(htmlPath)
	if err != nil {
		http.Error(w, "Error reading HTML file", http.StatusInternalServerError)
		return
	}

	// Replace placeholder in the HTML with the actual link
	output := bytes.Replace(htmlContent, []byte("{{LINK}}"), []byte(link), -1)

	w.Header().Set("Content-Type", "text/html")
	w.Write(output)
}

// Modify the submitHandler function to redirect to the new link display page
func SubmitHandler(w http.ResponseWriter, r *http.Request, connsaveshare *tarantool.Connection, prefix string) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		data := r.FormValue("data")
		if len(data) >= 10*1024*1024 {
			http.Error(w, "Too many symbols, maximum allowed size is 10MB", http.StatusBadRequest)
			return // end
		}
		hash := create_hash(strconv.FormatInt(time.Now().Unix(), 10))
		data_key, _ := create_random_key()
		data_json := put_data_in_json(hash, data)
		encrypted_data_json, _ := encrypt(data_json, data_key)

		if connsaveshare == nil {
			var err error
			connsaveshare, err = tarantool_connect_to_saveshare()
			if err != nil {
				http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
				return
			}
		}

		if err := add_new_data(encrypted_data_json, hash, connsaveshare); err != nil {
			fmt.Println(err)
			http.Error(w, "Failed to add new data", http.StatusInternalServerError)
			return
		}
		fmt.Println(prefix, hash, data_key)
		link := create_link(prefix, hash, data_key)
		// Redirect to the link display page with the generated link
		http.Redirect(w, r, "/link?link="+url.QueryEscape(link), http.StatusSeeOther)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func SaveshareRedirectPageHandler(w http.ResponseWriter, r *http.Request, path string, connsaveshare *tarantool.Connection) {
	// Extract the hash and key from the URL
	parts := strings.Split(r.URL.Path, "/ss/")
	if len(parts) != 2 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		fmt.Println("len!=2")
		return
	}

	// Split the hash and key
	hashKeyParts := strings.Split(parts[1], "&key=")
	if len(hashKeyParts) != 2 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		fmt.Println("hashkeyparts !=2")
		return
	}

	hash := hashKeyParts[0]
	keyHex := hashKeyParts[1]

	// Call the Tarantool function to get the data
	statusCode, dataTuple, err := get_data_from_tarantool(hash, connsaveshare)
	if err != nil || statusCode != "200" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		fmt.Println("statuscode != 200")
		return
	}

	// Decrypt the data
	encryptedData := dataTuple[1].(string) // Assuming the second element is the encrypted data
	dataKey, err := hex.DecodeString(keyHex)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		fmt.Println("key wrong")
		return
	}

	decryptedData, err := decrypt(encryptedData, dataKey)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		fmt.Println("unlucky decrypt", err)
		return
	}

	// Validate the decrypted data
	if !if_valid_json(decryptedData, hash) {
		fmt.Println("json is not valid")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	//fmt.Println(reflect.TypeOf(decryptedData)) json_data
	// Return the decrypted data here at ss.html
	data_to_print := decryptedData.Data                           //this is a data we need to xfer to user
	t, err := template.ParseFiles(filepath.Join(path, "ss.html")) //
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, data_to_print)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
