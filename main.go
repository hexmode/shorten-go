package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
)

var configMap map[string]string

func init() {
	configMap = make(map[string]string)
	configMap["title"] = "amd.im"
	configMap["length"] = "5"
}

var confMux = &sync.Mutex{}

func Config(name string) string {
	confMux.Lock()
	defer confMux.Unlock()
	return configMap[name]
}

type PageVars struct {
	Title string
	New   string
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "GET amd.im")
}

func NewHandler(w http.ResponseWriter, r *http.Request) {

	length, err := strconv.Atoi(Config("length"))
	if err != nil {
		log.Fatal(err)
	}

	newVars := PageVars{
		Title: "New Short URL | " + Config("title"),
		New:   generateKey(length),
	}

	t, err := template.ParseFiles("new.html")
	if err != nil {
		log.Fatal(err)
	}

	err = t.Execute(w, newVars)
	if err != nil {
		log.Fatal(err)
	}
}

func Redirector(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "GET amd.im/%v", vars["short"])
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/{short}", Redirector).Methods("GET")
	r.HandleFunc("/new", NewHandler).Methods("POST")
	r.HandleFunc("/", HomeHandler)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8000", r))
}
