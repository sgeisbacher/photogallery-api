package rest

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/gobuffalo/packr"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Server struct {
	RestLabelsHandler *RestLabelsHandler
	RestMediaHandler  *RestMediaHandler
	MediaFilesHandler http.Handler
}

var REGEXP_DATA_URL = regexp.MustCompile(`^/data/media/(big|thumb|orig)/(\w+)$`)

func (srv *Server) Serve(staticBox packr.Box) {
	router := mux.NewRouter().StrictSlash(true)
	router.Methods("GET").Path("/labels").HandlerFunc(srv.RestLabelsHandler.handleGetLabels)
	router.Methods("GET").Path("/labels/{id}").HandlerFunc(srv.RestLabelsHandler.handleGetLabel)
	router.Methods("GET").Path("/medias").HandlerFunc(srv.RestMediaHandler.handleGetMedias)
	router.Methods("GET").Path("/medias/{hash}").HandlerFunc(srv.RestMediaHandler.handleGetMedia)
	router.Methods("GET").PathPrefix("/data/media/").Handler(srv.MediaFilesHandler)
	router.PathPrefix("/").Handler(logger(http.FileServer(staticBox)))
	handler := cors.Default().Handler(router)

	address := ":8080"
	fmt.Printf("server listening on %v ...", address)
	log.Fatal(http.ListenAndServe(address, handler))
}

func logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("http:", r.RequestURI)
		h.ServeHTTP(w, r)
	})
}
