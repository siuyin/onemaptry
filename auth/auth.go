package auth

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

type token struct {
	Token     string `json:"access_token"`
	Timestamp int    `json:"timestamp"`
}

func Token() (string, error) {
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
		return "", fmt.Errorf("new request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := cl.Do(req)
	if err != nil {
		return "", fmt.Errorf("client do: %v", err)
	}

	defer resp.Body.Close()
	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("readall: ", err)
	}

	var tok token
	if err := json.Unmarshal(dat, &tok); err != nil {
		return "", fmt.Errorf("unmarshal: %v", err)
	}
	return tok.Token, nil
}

func creds() *bytes.Buffer {
	var buf bytes.Buffer
	je := json.NewEncoder(&buf)
	if err := je.Encode(auth{dflt.EnvString("EMAIL", "your email"), dflt.EnvString("PASSWD", "your passwd")}); err != nil {
		log.Fatal("encode: ", err)
	}
	return &buf
}

func Request(method string, url string, body io.ReadCloser) (*http.Response, error) {
	cl := &http.Client{}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("new request: %v", err)
	}

	req.Header.Add("Authorization", "Bearer "+dflt.EnvString("TOKEN", "mytoken"))
	resp, err := cl.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client do: %v", err)
	}

	return resp, nil

}
