package main

import (
	"fmt"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/sgeisbacher/photogallery-api/importer"
	"github.com/sgeisbacher/photogallery-api/media"
	"github.com/sgeisbacher/photogallery-api/rest"
)

func main() {
	fmt.Println("starting")
	db, err := bolt.Open("./data/data.db", 0600, nil)
	if err != nil {
		fmt.Println("error while opening db:", err)
	}
	galleryService := &media.GalleryService{db}
	mediaService := &media.MediaService{
		Db:             db,
		GalleryService: galleryService,
	}
	importManager := importer.ImportManager{
		MediaService:   mediaService,
		GalleryService: galleryService,
	}
	go importManager.ScanFolder("./data/orig")
	restServer := rest.Server{
		RestGalleryHandler: &rest.RestGalleryHandler{galleryService},
		RestMediaHandler:   &rest.RestMediaHandler{mediaService},
	}
	restServer.Serve()
	fmt.Println("done!")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	responseStr := `<h1>this is photogallery-api</h1><p>...</p>`
	fmt.Fprintln(w, responseStr)
}
