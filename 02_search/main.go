package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"onemaptry/auth"
	"onemaptry/body"
	"time"

	"github.com/siuyin/dflt"
)

const baseURL = "https://onemap.gov.sg"

type result struct {
	X   string
	Y   string
	Lat string `json:"LATITUDE"`
	Lng string `json:"LONGITUDE"`
}

func (r result) String() string {
	return fmt.Sprintf("(X,Y)=(%s,%s), LatLng=(%s,%s)", r.X, r.Y, r.Lat, r.Lng)
}

type resp struct {
	Found   int      `json:"found"`
	Results []result `json:"results"`
}

func main() {
	loc := url.QueryEscape(dflt.EnvString("LOC", "revenue house"))
	dat := retryOnUnauth(search, loc)
	fmt.Printf("%s\n", dat)

	var res resp
	if err := json.Unmarshal(dat, &res); err != nil {
		log.Fatal("unmarshall: ", err)
	}

	fmt.Printf("%s\n", res.Results)
}

func retryOnUnauth(fn func(string) ([]byte, error), loc string) []byte {
	done := false
	for dat, err := fn(loc); !done; dat, err = fn(loc) {
		if err != nil && err.Error() == "unauthorized" {
			_, err := auth.Token()
			if err != nil {
				log.Fatal("could not renew token: ", err)
			}

			log.Println("token refreshed")
			time.Sleep(1500 * time.Millisecond) // wait for new token to be registered in SLA's system
			continue
		}
		if err != nil {
			log.Fatal(err)
		}

		done = true
		return dat
	}
	return []byte{}
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
