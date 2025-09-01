package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/siuyin/dflt"
	"github.com/siuyin/gmap/lta/bike"
	"github.com/siuyin/onemaptry/05_defaultMap/public"
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
