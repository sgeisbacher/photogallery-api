package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sgeisbacher/photogallery-api/media"
)

type RestMediaHandler struct {
	MediaService *media.MediaService
}

func (handler *RestMediaHandler) handleGetMedias(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	medias, err := handler.MediaService.FindAll()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	if err := json.NewEncoder(w).Encode(medias); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
