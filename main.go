package main

import (
	"fmt"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/robfig/cron"
	"github.com/sgeisbacher/photogallery-api/imageconvertion"
	"github.com/sgeisbacher/photogallery-api/importer"
	"github.com/sgeisbacher/photogallery-api/media"
	"github.com/sgeisbacher/photogallery-api/metadata"
	"github.com/sgeisbacher/photogallery-api/rest"
)

func main() {
	fmt.Println("starting")

	// create db
	db := createDB()

	// create services
	galleryService := &media.GalleryService{db}
	mediaService := createMediaService(db, galleryService)
	importManager := createImportManager(mediaService, galleryService)
	metaDataManager := createMetaDataManager(mediaService)

	// start importer
	go importManager.ScanFolder("./data/orig")

	// set up cronjobs
	cronJobs := cron.New()
	cronJobs.AddFunc("@every 30s", func() { metaDataManager.Run() })
	cronJobs.Start()

	// create and start RestServer
	restServer := createRestServer(galleryService, mediaService)
	restServer.Serve()

	fmt.Println("done!")
}

func createDB() *bolt.DB {
	db, err := bolt.Open("./data/data.db", 0600, nil)
	if err != nil {
		panic(fmt.Sprintf("could not create db: %v", err))
	}
	return db
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	responseStr := `<h1>this is photogallery-api</h1><p>...</p>`
	fmt.Fprintln(w, responseStr)
}

func createMediaService(db *bolt.DB, galleryService *media.GalleryService) *media.MediaService {
	return &media.MediaService{
		Db:             db,
		GalleryService: galleryService,
	}
}

func createImportManager(mediaService *media.MediaService, galleryService *media.GalleryService) importer.ImportManager {
	return importer.ImportManager{
		MediaService:   mediaService,
		GalleryService: galleryService,
	}
}

func createRestServer(galleryService *media.GalleryService, mediaService *media.MediaService) rest.Server {
	hashFileSystem := rest.HashFileSystem{
		DataRoot:     "",
		MediaService: mediaService,
	}
	fileSystemLogDecorator := rest.FileSystemLogDecorator{
		FileSystem: hashFileSystem,
	}
	return rest.Server{
		RestGalleryHandler: &rest.RestGalleryHandler{galleryService},
		RestMediaHandler:   &rest.RestMediaHandler{mediaService},
		MediaFilesHandler:  http.FileServer(fileSystemLogDecorator),
	}
}

func createMetaDataManager(mediaService *media.MediaService) *metadata.MetaDataManager {
	var handlers []metadata.MetaDataHandler
	handlers = append(handlers, metadata.ShootTimeMetaDataHandler{})
	handlers = append(handlers, metadata.ThumbnailHandler{
		imageconvertion.ImageMagickImageConverter{},
	})
	return &metadata.MetaDataManager{
		MediaService:     mediaService,
		MetaDataHandlers: handlers,
	}
}
