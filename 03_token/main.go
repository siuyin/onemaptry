package main

import (
	"fmt"
	"log"
	"onemaptry/auth"
	"os"

	"github.com/siuyin/dflt"
)

func main() {
	tok, err := auth.Token()
	if err != nil {
		log.Fatal("auth.token: ", err)
	}

	os.Setenv("TOKEN", tok)
	fmt.Println(dflt.EnvString("TOKEN", "my token"))
}
