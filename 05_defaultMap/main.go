package main

import (
	"bytes"
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
	t = template.Must(template.ParseFS(public.Content, "tmpl/*"))

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

	t := template.Must(template.New("pe").Parse(
		`<div id="results">
		{{.Found}} result(s) found. page: {{.PageNum}} of {{.Pages}}
		<ul>
		{{range .Results}}
		  <li><a href="#" data-on-click="@get('/center?lat={{.Lat}}&lng={{.Lng}}&addr={{.Address}}')">{{.Address}}</a></li>
		{{end}}
		</ul>
		</div>`,
	))
	var b bytes.Buffer
	t.Execute(&b, sr)
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
	fmt.Printf("%s\n", r.URL)
	sse := datastar.NewSSE(w, r)
	sse.ExecuteScript(`markers.clearLayers()`)
	sse.ExecuteScript(fmt.Sprintf(`map.setView([%s, %s],18)`, r.FormValue("lat"), r.FormValue("lng")))
	sse.ExecuteScript(fmt.Sprintf(`markers.addLayer(L.marker([%s,%s]).bindPopup("%s"));markers.addTo(map)`, r.FormValue("lat"), r.FormValue("lng"), r.FormValue("addr")))
}
