package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Templates for each page
var idxtmpl *template.Template
var restmpl *template.Template
var uuidAPIResponses = make(map[string]string)

// QueryData used to generate a response with the LLM
type QueryData struct {
	URLcsv     string
	HumanQuery string
	UUID       string
}

type ResponsePageData struct {
	UUID string
}

// PageData to populate the page
type PageData struct {
	Title string
}

type ApiData struct {
	Resources string
	Query     string
}

// Index handling function
func index(w http.ResponseWriter, r *http.Request) {

	data := PageData{
		Title: "T.E.S.",
	}

	idxtmpl.Execute(w, data)
}

func response(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["UUID"]
	respgdata := ResponsePageData{
		UUID: uuid,
	}

	restmpl.Execute(w, respgdata)
}

// Query handling function
func query(w http.ResponseWriter, r *http.Request) {
	id := uuid.New()
	w.Header().Set("Content-Type", "application/json")

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	queryData := QueryData{
		URLcsv:     r.Form.Get("urlcsv"),
		HumanQuery: r.Form.Get("humanquery"),
		UUID:       id.String(),
	}
	//TODO actually go get url data and place it in a new struct
	fmt.Println(queryData.HumanQuery)
	fmt.Println(queryData.URLcsv)
	fmt.Println(queryData.UUID)

	sendData := url.Values{}
	sendData.Set("Resources", queryData.URLcsv)
	sendData.Set("Query", queryData.HumanQuery)
	sendData.Set("UUID", queryData.UUID)
	_, err = http.PostForm("http://localhost:9090/query", sendData)
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/response/"+id.String(), 303)

}
func sync(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["UUID"]
	if val, ok := uuidAPIResponses[uuid]; ok {
		response := map[string]string{"response": val}
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func serverResponse(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["UUID"]
	err := r.ParseForm()
	if err != nil {
		// Might want to do something here to tell the client their response might not be coming
		return
	}
	resp := r.Form.Get("Response")
	uuidAPIResponses[uuid] = resp
}

func main() {
	router := mux.NewRouter()
	idxtmpl = template.Must(template.ParseFiles("templates/index.html"))
	restmpl = template.Must(template.ParseFiles("templates/resp.html"))

	fs := http.FileServer(http.Dir("./static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/query", query).Methods("POST")
	router.HandleFunc("/response/{UUID}", response).Methods("GET")
	router.HandleFunc("/response/{UUID}/sync", sync).Methods("GET")
	router.HandleFunc("/server-response/{UUID}", serverResponse).Methods("POST")
	fmt.Println("Listening on :9091")
	log.Fatal(http.ListenAndServe(":9091", router))
}
