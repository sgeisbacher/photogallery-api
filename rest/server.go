package rest

import (
	"log"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Server struct {
	RestLabelsHandler *RestLabelsHandler
	RestMediaHandler  *RestMediaHandler
	MediaFilesHandler http.Handler
}

var REGEXP_DATA_URL = regexp.MustCompile(`^/data/media/(big|thumb|orig)/(\w+)$`)

func (srv *Server) Serve() {
	router := mux.NewRouter().StrictSlash(true)
	router.Methods("GET").Path("/labels").HandlerFunc(srv.RestLabelsHandler.handleGetLabels)
	router.Methods("GET").Path("/labels/{id}").HandlerFunc(srv.RestLabelsHandler.handleGetLabel)
	router.Methods("GET").Path("/medias").HandlerFunc(srv.RestMediaHandler.handleGetMedias)
	router.Methods("GET").Path("/medias/{hash}").HandlerFunc(srv.RestMediaHandler.handleGetMedia)
	router.Methods("GET").PathPrefix("/data/media/").Handler(srv.MediaFilesHandler)
	handler := cors.Default().Handler(router)

	log.Fatal(http.ListenAndServe("127.0.0.1:8080", handler))
}
