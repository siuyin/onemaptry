package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/siuyin/dflt"
)

const baseURL = "https://onemap.gov.sg"

func main() {
	cl := &http.Client{}

	loc := url.QueryEscape(dflt.EnvString("LOC", "revenue house"))
	req, err := http.NewRequest("GET", baseURL+"/api/common/elastic/search?searchVal="+loc+"&returnGeom=Y&getAddrDetails=Y&pageNum=1", nil)
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
