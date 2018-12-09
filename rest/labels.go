package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sgeisbacher/goutils/webutils"
	"github.com/sgeisbacher/photogallery-api/labels"
	"github.com/sgeisbacher/photogallery-api/media"
)

type RestLabelsHandler struct {
}

type LabelDTO struct {
	ID     string
	Name   string
	Medias []MediaDTO
}

func (handler *RestLabelsHandler) handleGetLabels(w http.ResponseWriter, req *http.Request) {
	labels := labels.FindAll()
	webutils.RespondWithJSON(w, labels, false)
}

func (handler *RestLabelsHandler) handleGetLabel(w http.ResponseWriter, req *http.Request) {
	id, ok := mux.Vars(req)["id"]
	var err error
	var label *labels.Label
	if ok {
		label, err = labels.Find(id)
	}
	var mediaDTOs []MediaDTO
	labelMedias, err := labels.GetMedias(label.ID)
	if err != nil {
		fmt.Printf("error: could not find medias for label %q: %v\n", label.ID, err)
		webutils.RespondWithJSON(w, nil, true)
		return
	}
	for _, mhash := range labelMedias {
		m, err := media.Find(mhash)
		if err != nil {
			fmt.Printf("warn: skipping media %q: %v\n", mhash, err)
			continue
		}
		mdto, err := CreateMediaDTO(m, nil)
		mediaDTOs = append(mediaDTOs, *mdto)
	}
	labelDTO := LabelDTO{
		ID:     label.ID,
		Name:   label.Name,
		Medias: mediaDTOs,
	}
	webutils.RespondWithJSON(w, labelDTO, err != nil)
}
