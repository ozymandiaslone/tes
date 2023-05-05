package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var idxtmpl *template.Template
var restmpl *template.Template
var uuidAPIResponses = make(map[string]string)
var servers []string

type QueryData struct {
	URLcsv     string
	HumanQuery string
	UUID       string
}

type ResponsePageData struct {
	UUID string
}

type PageData struct {
	Title string
}

type ApiData struct {
	Resources string
	Query     string
}

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

	fmt.Println(queryData.HumanQuery)
	fmt.Println(queryData.URLcsv)
	fmt.Println(queryData.UUID)

	sendData := url.Values{}
	sendData.Set("Resources", queryData.URLcsv)
	sendData.Set("Query", queryData.HumanQuery)
	sendData.Set("UUID", queryData.UUID)
	_, err = http.PostForm("http://localhost:9090/query", sendData)
	if err != nil {
		fmt.Println(err)
		//Do something about this
		return
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

func popIdx(servers *[]string, idx int) {
	*servers = append((*servers)[:idx], (*servers)[idx+1:]...)
}

func startServer(router *mux.Router) {
	log.Fatal(http.ListenAndServe(":9091", router))
}

func checkAlive(url string) bool {
	resp, err := http.Get(url + "/pulse")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		return true
	}
	return false
}

func manageServers(servers *[]string) {
	for i := 0; i < len(*servers); i++ {
		if !checkAlive((*servers)[i]) {
			popIdx(servers, i)
		}
	}
	time.Sleep(time.Second * 2)
}

func serverListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(servers)
}

func startupPingHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ip := r.Form.Get("ip")
	servers = append(servers, ip)
}

func main() {
	router := mux.NewRouter()
	tmplDir, err := filepath.Abs("templates")
	if err != nil {
		log.Fatal(err)
	}
	idxtmpl = template.Must(template.ParseFiles(filepath.Join(tmplDir, "index.html")))
	restmpl = template.Must(template.ParseFiles(filepath.Join(tmplDir, "resp.html")))
	fs := http.FileServer(http.Dir("./static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/query", query).Methods("POST")
	router.HandleFunc("/response/{UUID}", response).Methods("GET")
	router.HandleFunc("/response/{UUID}/sync", sync).Methods("GET")
	router.HandleFunc("/server-response/{UUID}", serverResponse).Methods("POST")
	router.HandleFunc("/startup-ping", startupPingHandler).Methods("POST")
	router.HandleFunc("/server-list", serverListHandler)
	go startServer(router)
	go manageServers(&servers)
	fmt.Println("Listening on :9091")
	select {}
}
