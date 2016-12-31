package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sgeisbacher/photogallery-api/importer"
)

func main() {
	fmt.Println("starting")
	importManager := importer.ImportManager{}
	go importManager.ScanFolder("./data/orig")
	http.HandleFunc("/", indexHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
	fmt.Println("done!")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	responseStr := `<h1>this is photogallery-api</h1><p>...</p>`
	fmt.Fprintln(w, responseStr)
}
