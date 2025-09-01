package main

import (
	"fmt"
	"log"
	"os"

	"github.com/siuyin/onemaptry/auth"
	"github.com/siuyin/onemaptry/body"
)

const baseURL = "https://onemap.gov.sg"

func main() {
	//fmt.Printf("%s\n", auth.RetryOnUnauth(themes))
	//fmt.Printf("%s\n", auth.RetryOnUnauth(retrieve, "family"))
	fmt.Printf("%s\n", auth.RetryOnUnauth(retrieve, "bicyclerack"))
}

func disp(fn func(p ...any) ([]byte, error), q ...any) {
	done := false
	for dat, err := fn(q...); !done; dat, err = fn(q...) {
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

func themes(p ...any) ([]byte, error) {
	url := baseURL + "/api/public/themesvc/getAllThemesInfo?moreInfo=Y"
	resp, err := auth.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("auth request: %v", err)
	}

	defer resp.Body.Close()
	return body.Read(resp)
}

func retrieve(p ...any) ([]byte, error) {
	if len(p) != 1 {
		return []byte{}, fmt.Errorf("retrieve: unexpected input format")
	}

	theme := p[0].(string)
	url := baseURL + "/api/public/themesvc/retrieveTheme?queryName=" + theme
	resp, err := auth.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("auth request: %v", err)
	}

	defer resp.Body.Close()
	return body.Read(resp)
}
