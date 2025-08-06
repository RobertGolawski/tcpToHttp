package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"tcpToHttp/internal/request"
	"tcpToHttp/internal/response"
	"tcpToHttp/internal/server"
)

const port = 42069

func handler(w *response.Writer, req *request.Request) {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		body := []byte(`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`)
		err := w.WriteStatusLine(response.StatusBadRequest)
		if err != nil {
			log.Printf("error writing status line: %v", err)
			return
		}
		h := response.GetDefaultHeaders(len(body))
		err = w.WriteHeaders(h)
		if err != nil {
			log.Printf("error writing headers: %v", err)
			return
		}
		err = w.WriteBody(body)
		if err != nil {
			log.Printf("error writing body: %v", err)
			return
		}
		return
	case "/myproblem":
		body := []byte(`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`)
		err := w.WriteStatusLine(response.StatusInternalServerError)
		if err != nil {
			log.Printf("error writing status line: %v", err)
			return
		}
		h := response.GetDefaultHeaders(len(body))
		err = w.WriteHeaders(h)
		if err != nil {
			log.Printf("error writing headers: %v", err)
			return
		}
		err = w.WriteBody(body)
		if err != nil {
			log.Printf("error writing body: %v", err)
			return
		}
		return
	default:
		body := []byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`)
		err := w.WriteStatusLine(response.StatusOK)
		if err != nil {
			log.Printf("error writing status line: %v", err)
			return
		}
		h := response.GetDefaultHeaders(len(body))
		err = w.WriteHeaders(h)
		if err != nil {
			log.Printf("error writing headers: %v", err)
			return
		}
		err = w.WriteBody(body)
		if err != nil {
			log.Printf("error writing body: %v", err)
			return
		}
	}
}

// func handler(w *response.Writer, req *request.Request) {
// 	err := w.WriteStatusLine(req.)
// }

func main() {
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
