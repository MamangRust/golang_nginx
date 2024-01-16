## Golang with nginx

In this tutorial, we will cover the setup of SSL for a Golang application served through Nginx. We'll generate self-signed SSL certificates for both Nginx and the Golang application.

## Prerequisites

Before you begin, make sure you have the following tools installed:

- Docker
- OpenSSL

## Generate SSL Certificates

### Nginx SSL Certificate

```sh
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout nginx.key -out nginx.crt

```

### Localhost SSL Certificate (for Golang)

```sh
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout key.pem -out cert.pem

```

## Docker-Compose Configuration

Here's a basic `docker-compose.yml` file for the Golang application, Nginx, and other services:

```yml
version: '3'
services:
  root:
    build:
      context: ./root
      dockerfile: Dockerfile
    ports:
      - '8080:8080'

  auth:
    build:
      context: ./auth
      dockerfile: Dockerfile
    ports:
      - '8081:8081'

  home:
    build:
      context: ./home
      dockerfile: Dockerfile
    ports:
      - '8082:8082'

  nginx:
    image: nginx:stable-alpine
    ports:
      - '80:80'
      - '443:443'
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf
      - ./nginx/nginx.crt:/etc/nginx/nginx.crt
      - ./nginx/nginx.key:/etc/nginx/nginx.key
      - ./ssl:/etc/nginx/ssl
```

## Nginx Configuration

Create an `nginx.conf` file inside the `nginx` directory:

```nginx

events {}

http {
  upstream go-app {
    least_conn;
    server root:8080 weight=10 max_fails=3 fail_timeout=30s;
    server auth:8081 weight=10 max_fails=3 fail_timeout=30s;
    server home:8082 weight=10 max_fails=3 fail_timeout=30s;
  }

  server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$host$request_uri;
  }

  server {
    listen 443 ssl;
    server_name your-domain.com;

    ssl_certificate /etc/nginx/nginx.crt;
    ssl_certificate_key /etc/nginx/nginx.key;

    location / {
      proxy_pass http://go-app;
      proxy_http_version 1.1;
      proxy_set_header Upgrade $http_upgrade;
      proxy_set_header Connection 'upgrade';
      proxy_set_header Host $host;
      proxy_cache_bypass $http_upgrade;
    }

    location ~ ^/auth {
      rewrite ^/auth/(.*) /$1 break;
      proxy_pass http://go-app;
      proxy_http_version 1.1;
      proxy_set_header Upgrade $http_upgrade;
      proxy_set_header Connection 'upgrade';
      proxy_set_header Host $host;
      proxy_cache_bypass $http_upgrade;
    }

    location ~ ^/home {
      rewrite ^/home/(.*) /$1 break;
      proxy_pass http://go-app;
      proxy_http_version 1.1;
      proxy_set_header Upgrade $http_upgrade;
      proxy_set_header Connection 'upgrade';
      proxy_set_header Host $host;
      proxy_cache_bypass $http_upgrade;
    }
  }
}

```

## Running the Setup

1. Generate SSL certificates as mentioned above.

2. Create the necessary Dockerfiles and Golang application code.

3. Place the docker-compose.yml, SSL certificates, and Nginx configuration in the appropriate directories.
4. Run docker-compose up -d to start the services.

Now, your Golang application should be accessible via HTTPS through Nginx. Visit https://your-domain.com in your browser to access the application securely.
