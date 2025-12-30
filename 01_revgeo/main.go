package main

import (
	"fmt"

	"github.com/siuyin/dflt"
	"github.com/siuyin/onemaptry/auth"
	"github.com/siuyin/onemaptry/body"
)

const baseURL = "https://onemap.gov.sg"

func main() {
	loc := dflt.EnvString("LOC", "1.3254295,103.9005321")
	fmt.Printf("%s", auth.RetryOnUnauth(reverse_geocode, loc))
}

func reverse_geocode(loc ...any) ([]byte, error) {
	url := baseURL + "/api/public/revgeocode?location=" + loc[0].(string) + "&buffer=20&addressType=All&otherFeatures=Y"
	resp, err := auth.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("auth.request: %v", err)
	}

	defer resp.Body.Close()
	return body.Read(resp)
}
