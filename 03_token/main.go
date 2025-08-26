package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/siuyin/dflt"
)

const baseURL = "https://onemap.gov.sg"

type auth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func main() {
	cl := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if via[len(via)-1].Method == "POST" && (req.Response.StatusCode == http.StatusFound) {
				req.Method = "POST"
				req.Body = io.NopCloser(creds())
				return nil
			}
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest("POST", baseURL+"/api/auth/post/getToken", creds())
	if err != nil {
		log.Fatal("new request: ", err)
	}

	req.Header.Add("Content-Type", "application/json")
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
func creds() *bytes.Buffer {
	var buf bytes.Buffer
	je := json.NewEncoder(&buf)
	if err := je.Encode(auth{dflt.EnvString("EMAIL", "your email"), dflt.EnvString("PASSWD", "your passwd")}); err != nil {
		log.Fatal("encode: ", err)
	}
	return &buf
}
