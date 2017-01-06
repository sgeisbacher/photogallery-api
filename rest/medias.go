package rest

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sgeisbacher/goutils/webutils"
	"github.com/sgeisbacher/photogallery-api/media"
)

type RestMediaHandler struct {
	MediaService *media.MediaService
}

func (handler *RestMediaHandler) handleGetMedias(w http.ResponseWriter, req *http.Request) {
	medias, err := handler.MediaService.FindAll()
	webutils.RespondWithJSON(w, medias, err != nil)
}

func (handler *RestMediaHandler) handleGetMedia(w http.ResponseWriter, req *http.Request) {
	hash, ok := mux.Vars(req)["hash"]
	var media *media.Media
	var err error
	if ok {
		media, err = handler.MediaService.FindMediaByHash(hash)
	}
	webutils.RespondWithJSON(w, media, err != nil)
}
