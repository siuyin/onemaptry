package main

import (
	"fmt"
	"log"
	"net/url"
	"onemaptry/auth"
	"onemaptry/body"
	"os"

	"github.com/siuyin/dflt"
)

const baseURL = "https://onemap.gov.sg"

func main() {
	loc := url.QueryEscape(dflt.EnvString("LOC", "revenue house"))
	done := false
	for dat, err := search(loc); !done; dat, err = search(loc) {
		if err != nil && err.Error() == "unauthorized" {
			tok, err := auth.Token()
			if err != nil {
				log.Fatal("could not renew token: ", err)
			}

			log.Println("token refreshed")
			os.Setenv("TOKEN", tok)
			continue
		}
		if err != nil {
			log.Fatal(err)
		}

		done = true
		fmt.Printf("%s\n", dat)
	}
}

func search(loc string) ([]byte, error) {
	url := baseURL + "/api/common/elastic/search?searchVal=" + loc + "&returnGeom=Y&getAddrDetails=Y&pageNum=1"
	resp, err := auth.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("auth request: %v", err)
	}

	defer resp.Body.Close()
	//fmt.Println(resp.Status)
	return body.Read(resp)
}
