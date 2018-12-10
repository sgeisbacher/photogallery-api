package rest

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sgeisbacher/goutils/webutils"
	"github.com/sgeisbacher/photogallery-api/media"
)

type RestMediaHandler struct {
}

type MediaDTO struct {
	Hash      string
	Name      string
	ThumbUrl  string
	BigUrl    string
	OrigUrl   string
	MediaType int
	ShootTime time.Time
	Labels    []string
}

func (handler *RestMediaHandler) handleGetMedias(w http.ResponseWriter, req *http.Request) {
	medias := media.FindAll()
	mediaDTOs, err := CreateMediaDTOs(medias)
	webutils.RespondWithJSON(w, mediaDTOs, err != nil)
}

func (handler *RestMediaHandler) handleGetMedia(w http.ResponseWriter, req *http.Request) {
	hash, ok := mux.Vars(req)["hash"]
	if !ok {
		webutils.RespondWithJSON(w, nil, true)
		return
	}
	m, err := media.Find(hash)
	if err != nil {
		webutils.RespondWithJSON(w, nil, true)
		return
	}
	labels, err := media.GetLabels(m)
	if err != nil {
		fmt.Printf("error while getting labels for media %q: %v\n", m.Hash, err)
		webutils.RespondWithJSON(w, nil, true)
		return
	}
	mediaDTO, err := CreateMediaDTO(m, labels)
	webutils.RespondWithJSON(w, mediaDTO, err != nil)
}

func CreateMediaDTOs(medias []*media.Media) ([]*MediaDTO, error) {
	var mediaDTOs []*MediaDTO
	for _, m := range medias {
		mdto, err := CreateMediaDTO(m, nil)
		if err != nil {
			return nil, err
		}
		mediaDTOs = append(mediaDTOs, mdto)
	}
	return mediaDTOs, nil
}

func CreateMediaDTO(m *media.Media, labels []string) (*MediaDTO, error) {
	if m.Hash == "" {
		return nil, errors.New("could not create MediaDTO from empty Media")
	}
	return &MediaDTO{
		Hash:      m.Hash,
		Name:      m.Name,
		ThumbUrl:  "/data/media/thumb/" + m.Hash,
		BigUrl:    "/data/media/big/" + m.Hash,
		OrigUrl:   "/data/media/orig/" + m.Hash,
		MediaType: m.MediaType,
		ShootTime: m.ShootTime,
		Labels:    labels,
	}, nil
}
