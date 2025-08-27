package body

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func Read(resp *http.Response) ([]byte, error) {
	if resp.StatusCode == http.StatusUnauthorized {
		return []byte{}, fmt.Errorf("unauthorized")
	}

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("readall: %v", err)
	}

	if bytes.Contains(dat, []byte("error")) {
		return []byte{}, fmt.Errorf("unauthorized")
	}

	return dat, nil
}
