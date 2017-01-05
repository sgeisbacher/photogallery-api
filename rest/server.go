package rest

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	RestGalleryHandler *RestGalleryHandler
	RestMediaHandler   *RestMediaHandler
}

func (srv *Server) Serve() {
	router := mux.NewRouter().StrictSlash(true)
	router.Methods("GET").Path("/galleries").HandlerFunc(srv.RestGalleryHandler.handleGetGalleries)
	router.Methods("GET").Path("/medias").HandlerFunc(srv.RestMediaHandler.handleGetMedias)
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", router))
}
