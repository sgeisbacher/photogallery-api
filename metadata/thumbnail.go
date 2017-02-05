package metadata

import (
	"github.com/sgeisbacher/goutils/fileutils"
	"github.com/sgeisbacher/photogallery-api/imageconvertion"
)

type ThumbnailHandler struct {
	ImageConverter imageconvertion.ImageConverter
}

func (handler ThumbnailHandler) GetName() string {
	return "ThumbnailHandler"
}

func (handler ThumbnailHandler) UpdateNeeded(ctx *MetaDataHandlerContext) bool {
	return ctx.media.ThumbnailPath == "" || !fileutils.FileExists(ctx.media.ThumbnailPath)
}

func (handler ThumbnailHandler) Handle(ctx *MetaDataHandlerContext) error {
	path, err := handler.ImageConverter.Convert(*ctx.media, imageconvertion.DEFAULT_THUMBNAIL_MEDIAFORMAT)
	if err != nil {
		return err
	}

	ctx.media.ThumbnailPath = path
	return nil
}
