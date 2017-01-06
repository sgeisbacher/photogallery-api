package rest

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sgeisbacher/goutils/webutils"
	"github.com/sgeisbacher/photogallery-api/media"
)

type RestGalleryHandler struct {
	GalleryService *media.GalleryService
}

func (handler *RestGalleryHandler) handleGetGalleries(w http.ResponseWriter, req *http.Request) {
	galleries, err := handler.GalleryService.FindAll()
	webutils.RespondWithJSON(w, galleries, err != nil)
}

func (handler *RestGalleryHandler) handleGetGallery(w http.ResponseWriter, req *http.Request) {
	id, ok := mux.Vars(req)["id"]
	var gallery *media.Gallery
	var err error
	if ok {
		gallery, err = handler.GalleryService.FindGalleryById(id)
	}
	webutils.RespondWithJSON(w, gallery, err != nil)
}
