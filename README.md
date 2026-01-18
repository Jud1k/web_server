# TCP to HTTP Server
A minimal HTTP/1.1 server built from the ground up using raw TCP sockets in Go. This project was developed as part of the ["From TCP to HTTP"](https://www.boot.dev/courses/learn-http-protocol-golang) course on boot.dev, focusing on understanding the fundamental protocols that power the web.

## Features
Pure TCP implementation: No net/http package for core HTTP handling

HTTP/1.1 compliant: Implements the HTTP specification with proper message parsing

Chunked transfer encoding: Supports streaming responses with Transfer-Encoding: chunked

Trailer headers: Implements HTTP trailers for post-response metadata

HTTP proxy: Can act as a proxy to external services

## Getting Started
Installation
```bash
# Clone the repository
git clone <repository-url>
cd tcptohttp

# Run the server
go run ./cmd/httpserver [port_number]
```

The port number argument is optional. If not provided, the server defaults to port 8000.

Once running, access the server at: http://localhost:8000

## Available Routes
```bash
/
```
Welcome message with a simple text response.
```bash
/html-wrong
```
HTML response with a 400 Bad Request status code and humorous error message.
```bash
/html-server
```
HTML response with a 500 Internal Server Error status code.
```bash
/html-ok
```
HTML response with a 200 OK status code and sample content.
```bash
/httpbin
```
Proxies requests to https://httpbin.org using chunked transfer encoding. This endpoint demonstrates streaming responses and trailer headers with X-Content-SHA256 and X-Content-Length.

Example: GET /httpbin/stream/5 streams 5 JSON responses from httpbin.org.
```bash
/video
```
Streams a sample video file with proper chunked encoding, simulating a video streaming endpoint.

Need create dir assets and put .mp4 file in.
