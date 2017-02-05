package rest

import (
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sgeisbacher/goutils/webutils"
	"github.com/sgeisbacher/photogallery-api/media"
)

type RestMediaHandler struct {
	MediaService *media.MediaService
}

type MediaDTO struct {
	Hash      string
	Name      string
	ThumbUrl  string
	BigUrl    string
	OrigUrl   string
	MediaType int
	ShootTime time.Time
}

func (handler *RestMediaHandler) handleGetMedias(w http.ResponseWriter, req *http.Request) {
	medias, err := handler.MediaService.FindAll()
	webutils.RespondWithJSON(w, medias, err != nil)
}

func (handler *RestMediaHandler) handleGetMedia(w http.ResponseWriter, req *http.Request) {
	hash, ok := mux.Vars(req)["hash"]
	if !ok {
		webutils.RespondWithJSON(w, nil, true)
		return
	}
	media, err := handler.MediaService.FindMediaByHash(hash)
	if err != nil {
		webutils.RespondWithJSON(w, nil, true)
		return
	}
	mediaDTO, err := CreateMediaDTO(media)
	webutils.RespondWithJSON(w, mediaDTO, err != nil)
}

func CreateMediaDTO(media *media.Media) (*MediaDTO, error) {
	if media.Hash == "" {
		return nil, errors.New("could not create MediaDTO from empty Media")
	}
	return &MediaDTO{
		Hash:      media.Hash,
		Name:      media.Name,
		ThumbUrl:  "/data/media/thumb/" + media.Hash,
		BigUrl:    "/data/media/big/" + media.Hash,
		OrigUrl:   "/data/media/orig/" + media.Hash,
		MediaType: media.MediaType,
		ShootTime: media.ShootTime,
	}, nil
}
