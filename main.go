package main

import (
	"encoding/json"
	//"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	//"strconv"
	//"time"

	"github.com/codegangsta/negroni"
	//"github.com/gorilla/context"
	"github.com/gorilla/mux"
	//"github.com/kr/pretty"
)

func main() {
	router := mux.NewRouter().StrictSlash(false)
	fbrouter := mux.NewRouter().StrictSlash(false)

	router.PathPrefix("/facebook").Handler(negroni.New(
		negroni.Wrap(fbrouter),
	))

	/*
		fbrouter.HandleFunc("/facebook/webhook", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("hub.verify_token") == "laughingbatman" {
				w.Write([]byte(r.URL.Query().Get("hub.challenge")))
				return
			}

			w.Write([]byte("Error, wrong validation token"))
		}).Methods("GET", "POST")
	*/

	fbrouter.HandleFunc("/facebook/webhook", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Webhook request!")

		var callback Callback
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		if err != nil {
			panic(err)
		}

		if err := r.Body.Close(); err != nil {
			panic(err)
		}

		if err := json.Unmarshal(body, &callback); err != nil {
			http.Error(w, "Invalid callback.", 400)
			return
		}

		for _, msg := range callback.Entry[0].Messaging {
			log.Println(msg.Message.Text)
		}

		w.Write([]byte(""))
	}).Methods("GET", "POST")

	log.Println("Server started on port 3001")

	n := negroni.Classic()
	n.UseHandler(router)
	http.ListenAndServe(":3001", router)
}

type Callback struct {
	Object string          `json:"object"`
	Entry  []CallbackEntry `json:"entry"`
}

type CallbackEntry struct {
	ID        string `json:"id"`
	Time      uint64 `json:"time"`
	Messaging []CallbackEntryMessage
}

type CallbackEntryMessage struct {
	Sender struct {
		ID string `json:"id"`
	} `json:"sender"`
	Recipient struct {
		ID string `json:"id"`
	} `json:"recipient"`
	Timestamp uint64 `json:"timestamp"`
	Message   struct {
		MID  string `json:"mid"`
		Seq  uint64 `json:"seq"`
		Text string `json:"text"`
	}
}
