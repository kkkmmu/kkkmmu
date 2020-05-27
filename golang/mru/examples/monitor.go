package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func QueryRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		r.ParseForm()
		v := strings.TrimSpace(r.FormValue("name"))
		if v == "" {
			io.WriteString(w, "You must gvim the register name to dump")
			return
		}
		log.Println(v)
		io.WriteString(w, fmt.Sprintf("Query Register %s", v))
	}
}

func QueryTable(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		r.ParseForm()
		v := strings.TrimSpace(r.FormValue("name"))
		if v == "" {
			io.WriteString(w, "You must gvim the table name to dump")
			return
		}
		log.Println(v)
		io.WriteString(w, fmt.Sprintf("Query table %s", v))
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/register", QueryRegister)
	mux.HandleFunc("/table", QueryTable)

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}
