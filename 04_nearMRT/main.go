package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"onemaptry/auth"
	"os"

	"github.com/siuyin/dflt"
)

const baseURL = "https://onemap.gov.sg"

func main() {
	lat := dflt.EnvString("LAT", "1.4044693603639506")
	lng := dflt.EnvString("LNG", "103.90083627504391")
	rad := dflt.EnvString("RAD", "1000")

	done := false
	for dat, err := search(lat, lng, rad); !done; dat, err = search(lat, lng, rad) {
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

func search(lat, lng, rad string) ([]byte, error) {
	url := fmt.Sprintf("%s/api/public/nearbysvc/getNearestMrtStops?latitude=%s&longitude=%s&radius_in_meters=%s", baseURL, lat, lng, rad)
	resp, err := auth.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("auth request: %v", err)
	}

	defer resp.Body.Close()
	//fmt.Println(resp.Status)
	if resp.StatusCode == http.StatusUnauthorized {
		return []byte{}, fmt.Errorf("unauthorized")
	}

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("readall: %v", err)
	}

	if bytes.Contains(dat, []byte("error")) {
		return []byte{}, fmt.Errorf("unauthorized")
	}

	return dat, nil
}
