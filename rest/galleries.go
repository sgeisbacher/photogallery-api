package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sgeisbacher/photogallery-api/media"
)

type RestGalleryHandler struct {
	GalleryService *media.GalleryService
}

func (handler *RestGalleryHandler) handleGetGalleries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	galleries, err := handler.GalleryService.FindAll()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	if err := json.NewEncoder(w).Encode(galleries); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
