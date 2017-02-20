package web

import (
	"fmt"
	"log"
	"net/http"
)

type Http struct {
	server http.Server
}

func (w *Http) Listen(port int) {
	mux := http.NewServeMux()
	handler := &MyHandler{}
	mux.Handle("/favicon.ico", http.NotFoundHandler())
	mux.Handle("/", handler)
	log.Print("[WEB] Launching WEB on port :", port)
	server := http.Server{Handler: mux, Addr: fmt.Sprint(":", port)}
	log.Fatal(server.ListenAndServe())
}
