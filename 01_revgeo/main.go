package main

import (
	"fmt"
	"log"
	"onemaptry/auth"
	"onemaptry/body"

	"github.com/siuyin/dflt"
)

const baseURL = "https://onemap.gov.sg"

func main() {
	loc := dflt.EnvString("LOC", "1.3254295,103.9005321")
	done := false
	for dat, err := reverse_geocode(loc); !done; dat, err = reverse_geocode(loc) {
		if err != nil && err.Error() == "unauthorized" {
			tok, err := auth.Token()
			if err != nil {
				log.Fatal("could not renew token: ", err)
			}

			log.Println("token refreshed")
			log.Println("new token: ", tok)
			continue
		}
		if err != nil {
			log.Fatal(err)
		}

		done = true
		fmt.Printf("%s\n", dat)
	}
}

func reverse_geocode(loc string) ([]byte, error) {
	url := baseURL + "/api/public/revgeocode?location=" + loc + "&buffer=20&addressType=All&otherFeatures=Y"
	resp, err := auth.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("auth.request: %v")
	}

	defer resp.Body.Close()
	return body.Read(resp)
}
