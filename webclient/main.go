package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Template to execute later
var idxtmpl *template.Template

// QueryData used to generate a response with the LLM
type QueryData struct {
	URLcsv     string
	HumanQuery string
}

// PageData to populate the page
type PageData struct {
	Title string
}

// Index handling function
func index(w http.ResponseWriter, r *http.Request) {

	data := PageData{
		Title: "T.E.S.",
	}

	idxtmpl.Execute(w, data)
}

// Query handling function
func query(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	queryData := QueryData{
		URLcsv:     r.Form.Get("urlcsv"),
		HumanQuery: r.Form.Get("humanquery"),
	}

	fmt.Println(queryData.HumanQuery)
	fmt.Println(queryData.URLcsv)

}

func main() {
	router := mux.NewRouter()
	idxtmpl = template.Must(template.ParseFiles("templates/index.html"))

	fs := http.FileServer(http.Dir("./static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	router.HandleFunc("/", index)
	router.HandleFunc("/query", query).Methods("POST")
	fmt.Println("Listening on :9091")
	log.Fatal(http.ListenAndServe(":9091", router))
}
