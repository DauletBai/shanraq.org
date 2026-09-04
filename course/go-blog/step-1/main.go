// Step 1 — after the lesson "Your first program in Go: a web server".
//
// The whole blog is one line of text so far. Everything that follows is added
// to this file: the server never gets rewritten, it only grows.
package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Сәлем! Бұл менің блогым."))
	})

	log.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
