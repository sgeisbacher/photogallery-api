package main

import (
	"fmt"
	"net/http"

	// "github.com/robfig/cron"
	"github.com/sgeisbacher/photogallery-api/imageconvertion"
	"github.com/sgeisbacher/photogallery-api/importer"
	"github.com/sgeisbacher/photogallery-api/metadata"
	"github.com/sgeisbacher/photogallery-api/rest"
)

func main() {
	fmt.Println("starting")

	// create services
	importManager := importer.ImportManager{}
	metaDataManager := createMetaDataManager()

	// start importer
	importManager.ScanFolder("./data/orig")
	metaDataManager.Run()

	// set up cronjobs
	// cronJobs := cron.New()
	// cronJobs.AddFunc("@every 30s", func() { metaDataManager.Run() })
	// cronJobs.Start()

	// create and start RestServer
	restServer := createRestServer()
	restServer.Serve()

	fmt.Println("done!")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	responseStr := `<h1>this is photogallery-api</h1><p>...</p>`
	fmt.Fprintln(w, responseStr)
}

func createRestServer() rest.Server {
	hashFileSystem := rest.HashFileSystem{
		DataRoot: ".",
	}
	fileSystemLogDecorator := rest.FileSystemLogDecorator{
		FileSystem: hashFileSystem,
	}
	return rest.Server{
		RestLabelsHandler: &rest.RestLabelsHandler{},
		RestMediaHandler:  &rest.RestMediaHandler{},
		MediaFilesHandler: http.FileServer(fileSystemLogDecorator),
	}
}

func createMetaDataManager() *metadata.MetaDataManager {
	var handlers []metadata.MetaDataHandler
	handlers = append(handlers, metadata.ShootTimeMetaDataHandler{})
	handlers = append(handlers, metadata.ThumbnailHandler{
		imageconvertion.ImageMagickImageConverter{},
	})
	handlers = append(handlers, metadata.BigPhotoHandler{
		imageconvertion.ImageMagickImageConverter{},
	})
	return &metadata.MetaDataManager{
		MetaDataHandlers: handlers,
	}
}
