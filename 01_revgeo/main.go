package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/siuyin/dflt"
)

const baseURL = "https://onemap.gov.sg"

func main() {
	cl := &http.Client{}

	req, err := http.NewRequest("GET", baseURL+"/api/public/revgeocode?location="+dflt.EnvString("LOC", "1.3254295,103.9005321")+"&buffer=20&addressType=All&otherFeatures=N", nil)
	if err != nil {
		log.Fatal("new request: ", err)
	}

	req.Header.Add("Authorization", "Bearer "+dflt.EnvString("TOKEN", "mytoken"))
	resp, err := cl.Do(req)
	if err != nil {
		log.Fatal("do: ", err)
	}

	defer resp.Body.Close()
	fmt.Printf("Status: %s\n", resp.Status)
	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("readall: ", err)
	}

	fmt.Printf("%s\n", dat)

}
