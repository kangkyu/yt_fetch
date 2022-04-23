package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	generate "github.com/kangkyu/yt_fetch/generate_csv"
)

func main() {
	port := "8080"
	if value, ok := os.LookupEnv("PORT"); ok {
		port = value
	}
	http.HandleFunc("/", pageHandler)
	http.HandleFunc("/fetches", fetchHandler)
	log.Println("Listen on localhost:"+port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	t, err := template.ParseFiles("page.html")
	if err != nil {
		http.Error(w, "file not found", 404)
		return
	}
	t.Execute(w, "")
}

func fetchHandler(w http.ResponseWriter, r *http.Request) {
	channelID := r.FormValue("uuid")
	// TODO: need validation channelID presence

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=videosof"+channelID+".csv")

	if err := generate.GenerateCSV(w, channelID); err != nil {
		fmt.Fprint(w, "could not generate CSV:\n")
		fmt.Fprint(w, err.Error())
		// TODO: error response should be 400 or 500 by cases, how do you tell?
		// http.Error(w, "internal error", 500)
	}
}
