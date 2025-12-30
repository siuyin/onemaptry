package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
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

// Token retrieve the same JWT if unexpired or a new JWT if expired.
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
	put(tok.Token)
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

// Get makes an authenticated request by adding Autorization: Bearer <token> header.
func Get(url string) (*http.Response, error) {
	cl := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %v", err)
	}

	tok := get()
	//log.Println("current token: ", tok)
	req.Header.Add("Authorization", "Bearer "+tok)
	resp, err := cl.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client do: %v", err)
	}

	return resp, nil
}

var db *bolt.DB

func init() {
	secDB := dflt.EnvString("SECRETS_DB", "/tmp/onemap.db.secret")
	initDB(secDB)
	log.Printf("auth package SECRETS_DB=%s", secDB)
}

func initDB(path string) {
	var err error
	db, err = bolt.Open(path, 0600, nil)
	if err != nil {
		log.Fatal("initDB: ", err)
	}
	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("jwt"))
		return nil
	})
}

func put(val string) {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("jwt"))

		err := b.Put([]byte("token"), []byte(val))
		if err != nil {
			return fmt.Errorf("put token: %v", err)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func get() string {
	tok := ""
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("jwt"))
		t := b.Get([]byte("token"))
		if t == nil {
			tok = ""
			return nil
		}
		tok = string(t)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return tok
}

func RetryOnUnauth(fn func(...any) ([]byte, error), args ...any) []byte {
	done := false
	for dat, err := fn(args...); !done; dat, err = fn(args...) {
		if err != nil && err.Error() == "unauthorized" {
			_, err := Token()
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
