package main

import (
	"fmt"
	"log"

	"github.com/siuyin/onemaptry/auth"
)

func main() {
	tok, err := auth.Token()
	if err != nil {
		log.Fatal("auth.token: ", err)
	}

	fmt.Println(tok)
}
