server {
    server_name example.com www.example.com;

    # Error page configuration
    error_page 404 /404.html;  # Specify the path to your 404.html
    location = /404.html {
        internal;  # Handler for 404
        root /path/to/html;  # Specify the path to the directory with 404.html
    }

    # Main application proxy
    location / {
        proxy_pass http://localhost:8080; # Port where your application runs
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Static files location
    location /static/ {
        alias /path/to/static/files/; # Specify the correct path to static files
    }

    # Serve JavaScript files with the correct MIME type
    location ~ \.js$ {
        root /path/to/js/files;  # Set the root to the directory where your JS files are located
        add_header Content-Type application/javascript;  # Set the correct MIME type for JS files
        try_files $uri =404;  # Serve the JS file or return 404
    }

    # Serve specific images
    location /photo.png {
        alias /path/to/images/photo.png;
    }

    location /sad.png {
        alias /path/to/images/sad.png;
    }

    # SSL configuration
    listen 443 ssl; # managed by Certbot
    ssl_certificate /etc/letsencrypt/live/example.com/fullchain.pem; # managed by Certbot
    ssl_certificate_key /etc/letsencrypt/live/example.com/privkey.pem; # managed by Certbot
    include /etc/letsencrypt/options-ssl-nginx.conf; # managed by Certbot
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem; # managed by Certbot
}

server {
    listen 80;  # Listen on HTTP
    server_name example.com www.example.com 147.45.108.78;

    # Redirect all HTTP requests to HTTPS
    return 301 https://$host$request_uri; # managed by Certbot
}

# Optional: Set client_max_body_size globally
client_max_body_size 512M;  # This can be set in the server block if needed
