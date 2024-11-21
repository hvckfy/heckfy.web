# Heckfy Web

Welcome to the Heckfy Web project! This repository contains the source code and configuration for the Heckfy website, which utilizes Nginx as a web server and is built with Go (Golang).

## Project Structure

The main folder `heckfyweb.site` contains the following directories and files:
```
heckfyweb.site/
├── html/ # Contains HTML files for the website
├── images/ # Contains images used in the website
├── server/ # Contains Go server application
│ ├── go.mod # Go module file
│ ├── go.sum # Go dependencies
│ ├── web.go # Main application file
│ └── internal/ # Internal packages
│   ├── rql/ # Request Limiter package
│   │ └── rqlconstantas.go # Constants for RQL package
│   ├── saveshare/ # Pastebin-like application package
│   │ └── saveshareconstantas.go # Constants for SaveShare package
│   ├── tinyfy/ # URL shortening application package
│   │ └── tinyfyconstantas.go # Constants for Tinyfy package
│   └── home/ # Home page request handlers
│   └── notfound_navigate/ # Handlers for not found/navigate pages
├── tarantool_heckfy_shard1/ # Tarantool Tinyfy shard configuration
├── tarantool_rql_shard1/ # Tarantool RQL shard configuration
├── tarantool_saveshare_shard1/ # Tarantool SaveShare shard configuration
└── heckfy.ru # Nginx configuration file
```

## Features

- **Request Limiter**: The `rql` package implements a request limiter to control the amount of information requested for various features.
- **SaveShare**: The `saveshare` package is a secure version of a pastebin service, allowing users to save and share text securely.
- **Tinyfy**: The `tinyfy` package provides functionality to shorten URLs, making them easier to share.
- **Home Page**: The `home` package handles requests for the main homepage of the website.
- **Not Found Navigation**: The `notfound_navigate` package manages the handling of 404 errors and navigation for pages that do not exist.

## Nginx Configuration

The Nginx configuration file (`heckfy.ru`) is set up to serve the application and manage SSL certificates. It includes configurations for error pages, static file serving, and proxying requests to the Go application.

## Systemd Service

The project includes a systemd service file to ensure the Go application runs automatically on startup. This service manages the Go application and keeps it running in the background.

## Installation

To set up the project locally, follow these steps:

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/heckfyweb.site.git
   cd heckfyweb.site

    Install Go dependencies:

    cd server
    go mod tidy

    Build the Go application:

    go build -o app web.go

    Configure Nginx and start the service:
        Ensure your Nginx configuration is set up correctly.
        Start the systemd service for the Go application.

Database Connection Constants

Inside the tinyfyconstantas.go, you can find constants for database connection:

package tinyfy

const (
	Address  = "TARANTOOL IP"  // Replace with your Tarantool IP
	User     = "USERNAME"       // Replace with your username
	Password = "PASSWORD"       // Replace with your password
)

This allows for smooth changing of logins for databases.
Paths and Connection Constants in web.go

In web.go, you will find:
```go
type paths_struct struct {
	HTML string // Exported field (starts with uppercase)
}

var paths = paths_struct{HTML: "/path/to/htmls"} // Specify the path to your HTML files
var connsaveshare *tarantool.Connection

const prefix string = "https://example.com/" // Set your domain prefix
```
Usage

Once everything is set up, you can access the website at http://yourdomain.com or https://yourdomain.com.
Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue for any improvements or suggestions.
License

This project is licensed under the MIT License - see the LICENSE file for details.


### Notes:
- Replace `yourusername` in the clone command with your actual GitHub username.
- Adjust any paths and constants to reflect your actual setup as needed.
