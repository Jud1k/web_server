package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/Jud1k/web_server/internal/headers"
	"github.com/Jud1k/web_server/internal/request"
	"github.com/Jud1k/web_server/internal/response"
	"github.com/Jud1k/web_server/internal/server"
)

func handler(w *response.Writer, req *request.Request) {
	log.Println(string(req.Body))
	h := response.GetDefaultHeaders(0)
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/") {
		data := []byte(`Hello, its my implementation HTTP-server from The Mister "I worked in Netflix btw" in boot.dev.
I realy enjoy do this project. I think go is a awesome language and everyone should try it.
so if you read this and do not try programming in go, GO do it.`)
		h.Set("Content-Length", fmt.Sprint(len(data)))
		h.Set("Content-Type", "text/plain")
		w.WriteStatusLine(400)
		w.WriteHeaders(h)
		w.WriteBody(data)
	}
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/html-wrong") {
		data := []byte(`<html><head><title>400 Bad Request</title></head><body><h1>Bad Request</h1><p>Your request honestly kinda sucked.</p></body></html>`)
		h.Set("Content-Length", fmt.Sprint(len(data)))
		h.Set("Content-Type", "text/html")
		w.WriteStatusLine(400)
		w.WriteHeaders(h)
		w.WriteBody(data)
	}
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/html-server") {
		data := []byte(`<html><head><title>500 Internal Server Error</title></head><body><h1>Internal Server Error</h1><p>Okay, you know what? This one is on me.</p></body></html>`)
		h.Set("Content-Length", fmt.Sprint(len(data)))
		h.Set("Content-Type", "text/html")
		w.WriteStatusLine(500)
		w.WriteHeaders(h)
		w.WriteBody(data)
	}
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/html-ok") {
		data := []byte(`<html><head><title>200 OK</title></head><body><h1>Success!</h1><p>Your request was an absolute banger.</p></body></html>`)
		h.Set("Content-Length", fmt.Sprint(len(data)))
		h.Set("Content-Type", "text/html")
		w.WriteStatusLine(200)
		w.WriteHeaders(h)
		w.WriteBody(data)
	}
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {
		newTarget := fmt.Sprintf("https://httpbin.org%s", strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin"))
		resp, err := http.Get(newTarget)
		if err != nil {
			log.Printf("error: %s", err)
			return
		}

		w.WriteStatusLine(200)

		h.Del("Content-Length")
		h.Set("Transfer-Encoding", "chunked")
		h.Set("Trailer", "X-Content-SHA256, X-Content-Length")
		trailers := headers.Headers{}
		w.WriteHeaders(h)
		fullBody := []byte{}
		buf := make([]byte, 1024)
		for {
			n, err := resp.Body.Read(buf)
			if n > 0 {
				_, writeErr := w.WriteChunkedBody(buf[:n])
				fullBody = append(fullBody, buf[:n]...)
				if writeErr != nil {
					break
				}
			}
			if err != nil {
				if err != io.EOF {
					log.Printf("Read error: %s", err)
				}
				break
			}
		}
		w.WriteChunkedBodyDone()
		bodyHash := sha256.Sum256(fullBody)
		trailers.Set("X-Content-SHA256", fmt.Sprintf("%x", bodyHash))
		trailers.Set("X-Content-Length", fmt.Sprint(len(fullBody)))
		w.WriteTrailers(trailers)
	}
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/video") {
		fileName := "./assets/vim.mp4"
		video, err := os.ReadFile(fileName)
		if err != nil {
			log.Fatalf("error: cannot read file with name %s", fileName)
		}
		h.Set("Content-Type", "video/mp4")
		h.Set("Content-Length", fmt.Sprint(len(video)))
		w.WriteStatusLine(200)
		w.WriteHeaders(h)
		w.WriteBody(video)
	}
}

func main() {
	port := 8000
	if len(os.Args) > 1 {
		arg, err := strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		port = arg
	}
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
