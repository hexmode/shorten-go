package main

import (
	"fmt"
	"github.com/gorilla/mux"
	bolt "go.etcd.io/bbolt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

var db *bolt.DB

var configMap map[string]string

func init() {
	configMap = make(map[string]string)
	if os.Getenv("dbpath") != "" {
		configMap["dbpath"] = os.Getenv("dbpath")
	} else {
		configMap["dbpath"] = "./bolt.db"
	}
	configMap["title"] = "amd.im"
	configMap["domain"] = "amd.im"
	configMap["length"] = "5"
}

var confMux = &sync.Mutex{}

func Config(name string) string {
	confMux.Lock()
	defer confMux.Unlock()
	return configMap[name]
}

type PageVars struct {
	Title  string
	New    string
	Domain string
	Key    string
	URL    string
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	length, err := strconv.Atoi(Config("length"))
	if err != nil {
		log.Print(err)
	}

	newVars := PageVars{
		Title: Config("title"),
		New:   generateKey(length),
	}

	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal(err)
	}

	err = t.Execute(w, newVars)
	if err != nil {
		log.Fatal(err)
	}
}

func NewPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var rec Record
	if r.Form.Get("key") != "" {
		rec.Key = r.Form.Get("key")
	} else {
		length, err := strconv.Atoi(Config("length"))
		if err != nil {
			log.Print(err)
		}
		rec.Key = generateKey(length)
	}
	rec.Type = "URL"
	rec.URL = r.Form.Get("URL")

	err := saveRecord(rec)
	if err != nil {
		log.Print("Could not save record", rec.Key, err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "FAILED\nPOST amd.im/new\n", r.Form)
	}

	log.Printf("Created %v pointing to %v", rec.Key, rec.URL)

	newVars := PageVars{
		Title:  Config("title"),
		Domain: Config("domain"),
		Key:    rec.Key,
		URL:    rec.URL,
	}

	t, err := template.ParseFiles("templates/new.html")
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

	rec, err := getKey(vars["key"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 GET amd.im/key %v", vars["key"])
		return
	}
	log.Printf("Redirecting %v to %v", vars["key"], rec.URL)
	http.Redirect(w, r, rec.URL, http.StatusFound)
}

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func main() {

	var err error

	db, err = bolt.Open(Config("dbpath"), 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/new", NewPostHandler).Methods("POST")
	r.HandleFunc("/static/{file}", StaticHandler).Methods("GET")
	r.HandleFunc("/{key}", Redirector).Methods("GET")
	r.HandleFunc("/", HomeHandler)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8000", r))
}
