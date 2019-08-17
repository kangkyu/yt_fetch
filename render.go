package main

import (
	"html/template"
	"log"
	"net/http"
	"fmt"
)

func main() {
    http.HandleFunc("/", pageHandler)
    log.Println("Listen on localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
    t, err := template.ParseFiles("page.html")
    if err != nil {
    	fmt.Println("hello")
    }
    t.Execute(w, "")
}
