package main

import (
	//"encoding/json"
	//"fmt"
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
	fbrouter := mux.NewRouter()

	router.PathPrefix("/facebook").Handler(negroni.New(
		negroni.Wrap(fbrouter),
	))

	fbrouter.HandleFunc("/facebook/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("hub.verify_token") == "laughingbatman" {
			w.Write([]byte(r.URL.Query().Get("hub.challenge")))
		}

		w.Write([]byte("Error, wrong validation token"))
	})

	log.Println("Server started on port 3001")

	n := negroni.Classic()
	n.UseHandler(router)
	http.ListenAndServe(":3001", router)
}
