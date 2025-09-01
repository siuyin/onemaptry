package srch

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/siuyin/onemaptry/auth"
	"github.com/siuyin/onemaptry/body"
)

type Response struct {
	Found   int      `json:"found"`
	Pages   int      `json:"totalNumPages"`
	PageNum int      `json:"pageNum"`
	Results []Result `json:"results"`
}

type Result struct {
	SearchVal  string `json:"SEARCHVAL"`
	BlkNo      string `json:"BLK_NO"`
	RoadName   string `json:"ROAD_NAME"`
	Building   string `json:"BUILDING"`
	Address    string `json:"ADDRESS"`
	PostalCode string `json:"POSTAL"`
	X          string
	Y          string
	Lat        string `json:"LATITUDE"`
	Lng        string `json:"LONGITUDE"`
}

func (r Result) String() string {
	return fmt.Sprintf("%s: (X,Y)=(%s,%s), LatLng=(%s,%s)\n", r.Address, r.X, r.Y, r.Lat, r.Lng)
}

func Location(loc string) *Response {
	dat := auth.RetryOnUnauth(search, url.QueryEscape(loc))
	var r Response
	if err := json.Unmarshal(dat, &r); err != nil {
		log.Printf("ERROR: unmarshal: %s: %v", dat, err)
	}

	return &r
}

const baseURL = "https://onemap.gov.sg"

func search(loc ...any) ([]byte, error) {
	url := baseURL + "/api/common/elastic/search?searchVal=" + loc[0].(string) + "&returnGeom=Y&getAddrDetails=Y&pageNum=1"
	resp, err := auth.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("auth request: %v", err)
	}

	defer resp.Body.Close()
	//fmt.Println(resp.Status)
	return body.Read(resp)
}
