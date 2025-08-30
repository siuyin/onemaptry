package main

import (
	"fmt"
	"onemaptry/auth"
	"onemaptry/body"

	"github.com/siuyin/dflt"
)

const baseURL = "https://onemap.gov.sg"

func main() {
	lat := dflt.EnvString("LAT", "1.4044693603639506")
	lng := dflt.EnvString("LNG", "103.90083627504391")
	rad := dflt.EnvString("RAD", "1000")
	fmt.Printf("%s\n", auth.RetryOnUnauth(search, lat, lng, rad))
}

func search(args ...any) ([]byte, error) {
	url := fmt.Sprintf("%s/api/public/nearbysvc/getNearestMrtStops?latitude=%s&longitude=%s&radius_in_meters=%s", baseURL, args[0].(string), args[1].(string), args[2].(string))
	resp, err := auth.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("auth request: %v", err)
	}

	defer resp.Body.Close()
	//fmt.Println(resp.Status)
	return body.Read(resp)
}
