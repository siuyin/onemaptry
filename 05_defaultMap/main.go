package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/siuyin/dflt"
	"github.com/siuyin/gmap/lta/bike"
	"github.com/siuyin/onemaptry/05_defaultMap/public"
	"github.com/siuyin/onemaptry/srch"
	"github.com/starfederation/datastar-go/datastar"
)

var t *template.Template

func main() {
	port := dflt.EnvString("PORT", "8080")
	log.Printf("PORT=%s", port)
	t = template.New("mytpl").Funcs(template.FuncMap{"json": jsonify})
	t = template.Must(t.ParseFS(public.Content, "tmpl/*"))

	http.Handle("/{$}", http.HandlerFunc(indexHandler))
	//http.Handle("/", http.FileServer(http.Dir("./public")))
	http.Handle("/", http.FileServer(http.FS(public.Content)))
	http.HandleFunc("/placepicker", placePickerHandler)
	http.HandleFunc("/bicyclepark", bicyleParkHandler)
	http.HandleFunc("/bicycleRacks", bicycleRacksHandler)
	http.HandleFunc("/search", placeSearchHandler)
	http.HandleFunc("/center", centerHandler)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t.ExecuteTemplate(w, "index.html", nil)
}

func placePickerHandler(w http.ResponseWriter, r *http.Request) {
	t.ExecuteTemplate(w, "placepicker.html", nil)
}

func bicyleParkHandler(w http.ResponseWriter, r *http.Request) {
	//t.ExecuteTemplate(w, "bicyclepark.html", struct{ Key string }{key})
}

func bicycleRacksHandler(w http.ResponseWriter, r *http.Request) {
	lat, err := strconv.ParseFloat(r.FormValue("Lat"), 64)
	if err != nil {
		log.Fatal(err)
	}

	lng, err := strconv.ParseFloat(r.FormValue("Lng"), 64)
	if err != nil {
		log.Fatal(err)
	}

	io.WriteString(w, bike.ParkingSpots(lat, lng))
}

type placeInp struct {
	Search string `json:"search"`
}

func placeSearchHandler(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)
	pl := placeInp{}
	datastar.ReadSignals(r, &pl)
	if len(pl.Search) < 4 {
		sse.PatchElements(`<div id="results"></div>`)
		return
	}

	sr := srch.Location(pl.Search)
	if sr.Found == 0 {
		sse.PatchElements(`<div id="results">no results found</div>`)
		return
	}

	tm := template.New("pe").Funcs(template.FuncMap{"json": jsonify})
	tm = template.Must(tm.Parse(
		`<div id="results">
		  {{.Found}} result(s) found. page: {{.PageNum}} of {{.Pages}}
		  <ul>
		    {{range .Results}}
		      <li><a href="#" data-on-click="@get('/center?selected={{json .}}')">{{.Address}}</a></li>
		    {{end}}
		  </ul>
		</div>`,
	))
	var b bytes.Buffer
	tm.Execute(&b, sr)
	newFeat := `{
              "type": "FeatureCollection",
              "features": [
                      {
                              "type": "Feature",
                              "properties": {
                                      "name": "Woh Hup",
                                    },
                              "geometry": {
                                      "coordinates": [
                                              103.77258,
                                              1.345618
                                            ],
                                      "type": "Point"
                                    },
                              "id": 0
                            }
                    ]
            }
`
	sse.PatchElements(b.String())
	sse.ExecuteScript(fmt.Sprintf(`render(%s)`, newFeat))

}

func centerHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s\n", r.FormValue("selected"))
	s := selectedSearchResult(r)

	sse := datastar.NewSSE(w, r)
	sse.ExecuteScript(`markers.clearLayers()`)
	sse.ExecuteScript(fmt.Sprintf(`map.setView([%s, %s],18)`, s.Lat, s.Lng))
	sse.ExecuteScript(fmt.Sprintf(`markers.addLayer(L.marker([%s,%s]).bindPopup("%s"))`, s.Lat, s.Lng, s.Address))
	sse.ExecuteScript(fmt.Sprintf(`render(%s)`, bikeRacks(s.Lat, s.Lng)))
}

func bikeRacks(latitude, longitude string) string {
	lat, err := strconv.ParseFloat(latitude, 64)
	if err != nil {
		log.Printf("parse float: %s: %v", latitude, err)
	}

	lng, err := strconv.ParseFloat(longitude, 64)
	if err != nil {
		log.Printf("parse float: %s: %v", longitude, err)
	}

	return bike.ParkingSpots(lat, lng)
}

func selectedSearchResult(r *http.Request) *srch.Result {
	var s srch.Result
	if err := json.Unmarshal([]byte(r.FormValue("selected")), &s); err != nil {
		log.Printf("unmarshal: %s", r.FormValue("selected"), err)
	}
	return &s
}

func jsonify(v any) (template.JS, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("jsonify: %v", err)
	}
	return template.JS(data), nil
}
