package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	llama "github.com/go-skynet/go-llama.cpp"
	"github.com/gorilla/mux"
)

var (
	threads = 8
	tokens  = 128
	model   = "./server/models/wizardLM-7B.GGML.q4_2.bin"
	queue   []Queuer
	loadnum int
)

type IncomingData struct {
	Resources string
	Query     string
	UUID      string
}

type Queuer struct {
	Prompt string
	UUID   string
}

type ReturnData struct {
	Response string
	UUID     string
}

func popQueue(queue *[]Queuer) {
	*queue = (*queue)[1:]
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	in := new(IncomingData)
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}
	in.Query = r.FormValue("Query")
	in.Resources = r.FormValue("Resources")
	in.UUID = r.FormValue("UUID")

	prompt := in.Resources + "\n\n" + in.Query
	queue = append(queue, Queuer{prompt, in.UUID})
	loadnum += 1

}

func responseSender(uuid string, response string) {
	sendData := url.Values{
		"Response": {response},
	}
	resp, err := http.PostForm("http://localhost:9091/server-response/"+uuid, sendData)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
}

func queryKiller(job Queuer) {

	l, err := llama.New(model, llama.SetContext(128), llama.SetParts(-1))
	if err != nil {
		fmt.Println("Loading the model failed:", err.Error())
		os.Exit(1)
	}

	l.SetTokenCallback(nil)

	result, err := l.Predict(job.Prompt, llama.Debug, llama.SetTokens(tokens), llama.SetThreads(threads), llama.SetTopK(90), llama.SetTopP(0.86), llama.SetStopWords("llama"))
	if err != nil {
		panic(err)
	}

	responseSender(job.UUID, result)

	loadnum -= 1

}

func loadHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.Itoa(loadnum)))
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/query", queryHandler).Methods("POST")
	router.HandleFunc("load", loadHandler)
	fmt.Println("Starting server...")
	go func() {
		log.Fatal(http.ListenAndServe(":9090", router))
	}()
	fmt.Println("Listening on port :9090")

	// start the infinite loop in a goroutine
	go func() {
		for {
			if len(queue) != 0 {
				fmt.Println("Starting job...")
				queryKiller(queue[0])
				fmt.Println("Job complete. Popping queue.")
				popQueue(&queue)
				fmt.Println("Queue popped.")
				fmt.Println("Load number:", loadnum)
			}
		}
	}()

	// block the main function from exiting
	select {}
}
