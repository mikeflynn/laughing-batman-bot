package main

import (
	"encoding/json"
	//"fmt"
	"bytes"
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
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		if err := r.Body.Close(); err != nil {
			panic(err)
		}

		log.Println(string(body))

		if err := json.Unmarshal(body, &callback); err != nil {
			log.Println("Invalid callback")
			http.Error(w, "Invalid callback.", 400)
			return
		}

		for _, msg := range callback.Entry[0].Messaging {
			log.Println(msg.Message.Text)

			sendMessage(msg.Sender.ID, "Thanks for the message!")
		}

		w.Write([]byte(""))
	}).Methods("GET", "POST")

	log.Println("Server started on port 3001")

	n := negroni.Classic()
	n.UseHandler(router)
	http.ListenAndServe(":3001", router)
}

func sendMessage(userID string, text string) {
	msg := OutgoingMessage{}
	msg.Recipient.ID = userID
	msg.Message.Text = text
	jsonBytes, _ := json.Marshal(msg)

	req, err := http.NewRequest("POST", "https://graph.facebook.com/v2.6/me/messages", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	log.Println("Send Response Status:", resp.Status)
}

type Callback struct {
	Object string          `json:"object"`
	Entry  []CallbackEntry `json:"entry"`
}

type CallbackEntry struct {
	ID        string                 `json:"id"`
	Time      uint64                 `json:"time"`
	Messaging []CallbackEntryMessage `json:"messaging"`
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
	Attachments []struct {
		Type    string `json:"type"`
		Payload struct {
			URL string `json:"url"`
		} `json:"payload"`
	} `json:"attachments"`
	Delivery []struct {
		MIDs      []string `json:"mids"`
		Watermark uint64   `json:"watermark"`
		Seq       uint64   `json:"seq"`
	}
	Postback struct {
		Payload string `json:"payload"`
	}
}

type OutgoingMessage struct {
	Recipient struct {
		ID string `json:"id"`
	} `json:"recipient"`
	Message struct {
		Text string `json:"text"`
	} `json:"message"`
}
