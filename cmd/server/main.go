package main

import (
	"log"

	"github.com/drakelthedragon/proglog/internal/server"
)

func main() {
	srv := server.NewHTTPServer(":4000")
	log.Fatal(srv.ListenAndServe())
}
