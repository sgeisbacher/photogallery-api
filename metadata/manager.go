package metadata

import (
	"fmt"
	"os"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/sgeisbacher/photogallery-api/media"
)

type MetaDataHandler interface {
	GetName() string
	UpdateNeeded(ctx *MetaDataHandlerContext) bool
	Handle(ctx *MetaDataHandlerContext) error
}

type MetaDataManager struct {
	MetaDataHandlers []MetaDataHandler
}

type MetaDataHandlerContext struct {
	exifData *exif.Exif
	media    *media.Media
}

func (mgr *MetaDataManager) Run() {
	fmt.Println("starting MetaDataManager ...")
	if len(mgr.MetaDataHandlers) == 0 {
		fmt.Println("MetaDataManager: no handlers configured, skipping run ...")
		return
	}
	medias := media.FindAll()
	for _, m := range medias {
		changed := false
		ctx := &MetaDataHandlerContext{media: m}
		for _, handler := range mgr.MetaDataHandlers {
			if handler.UpdateNeeded(ctx) {
				fmt.Printf("processing '%v' on media '%v' ...\n", handler.GetName(), m.Hash)
				err := runHandler(handler, ctx)
				if err != nil {
					fmt.Print("error while processing '%v' on '%v': %v\n", handler.GetName(), m.Hash, err)
					continue
				}
				changed = true
			}
		}
		if changed {
			media.Add(m)
		}
	}
	fmt.Println("MetaDataManager ... done!")
}

func runHandler(handler MetaDataHandler, ctx *MetaDataHandlerContext) error {
	var err error
	if ctx.exifData == nil {
		ctx.exifData, err = extractExifData(ctx.media)
		if err != nil {
			fmt.Printf("error while loading exifdata for '%v': %v\n", ctx.media.OrigPath, err)
			return err
		}
	}

	return handler.Handle(ctx)
}

func extractExifData(media *media.Media) (*exif.Exif, error) {
	file, err := os.Open(media.OrigPath)
	if err != nil {
		return nil, err
	}

	return exif.Decode(file)
}
