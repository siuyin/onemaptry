package main

import (
	"fmt"
	"log"
	"onemaptry/auth"
)

func main() {
	tok, err := auth.Token()
	if err != nil {
		log.Fatal("auth.token: ", err)
	}

	fmt.Println(tok)
}
