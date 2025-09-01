package main

import (
	"fmt"

	"github.com/siuyin/dflt"
	"github.com/siuyin/onemaptry/auth"
	"github.com/siuyin/onemaptry/body"
)

const baseURL = "https://onemap.gov.sg"

func main() {
	lat := dflt.EnvString("LAT", "1.4044693603639506")
	lng := dflt.EnvString("LNG", "103.90083627504391")
	rad := dflt.EnvString("RAD", "1000")
	fmt.Printf("Nearest MRT: %s\n\n", auth.RetryOnUnauth(search, "Mrt", lat, lng, rad))
	fmt.Printf("Nearest Bus: %s\n", auth.RetryOnUnauth(search, "Bus", lat, lng, rad))
}

func search(args ...any) ([]byte, error) {
	url := fmt.Sprintf("%s/api/public/nearbysvc/getNearest%sStops?latitude=%s&longitude=%s&radius_in_meters=%s", baseURL, args[0].(string), args[1].(string), args[2].(string), args[3].(string))
	resp, err := auth.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("auth request: %v", err)
	}

	defer resp.Body.Close()
	//fmt.Println(resp.Status)
	return body.Read(resp)
}
