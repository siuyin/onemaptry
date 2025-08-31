package main

import (
	"log"
	"net/http"

	"github.com/siuyin/dflt"
)

func main() {
	port := dflt.EnvString("PORT", "8080")
	log.Printf("PORT=%s", port)
	http.Handle("/", http.FileServer(http.Dir("./public")))
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
