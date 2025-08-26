package main

import (
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
	loc := dflt.EnvString("LOC", "1.3254295,103.9005321")
	done := false
	for dat, err := reverse_geocode(loc); !done; dat, err = reverse_geocode(loc) {
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

func reverse_geocode(loc string) ([]byte, error) {
	cl := &http.Client{}

	req, err := http.NewRequest("GET", baseURL+"/api/public/revgeocode?location="+loc+"&buffer=20&addressType=All&otherFeatures=N", nil)
	if err != nil {
		return []byte{}, fmt.Errorf("new request: %v", err)
	}

	req.Header.Add("Authorization", "Bearer "+dflt.EnvString("TOKEN", "mytoken"))
	resp, err := cl.Do(req)
	if err != nil {
		return []byte{}, fmt.Errorf("client do: %v", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized {
		return []byte{}, fmt.Errorf("unauthorized")
	}

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("readall: %v", err)
	}

	return dat, nil
}
