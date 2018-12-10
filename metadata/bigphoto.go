package metadata

import (
	"github.com/sgeisbacher/goutils/fileutils"
	"github.com/sgeisbacher/photogallery-api/imageconvertion"
)

type BigPhotoHandler struct {
	ImageConverter imageconvertion.ImageConverter
}

func (handler BigPhotoHandler) GetName() string {
	return "BigPhotoHandler"
}

func (handler BigPhotoHandler) UpdateNeeded(ctx *MetaDataHandlerContext) bool {
	return ctx.media.Path == "" || !fileutils.FileExists(ctx.media.Path)
}

func (handler BigPhotoHandler) Handle(ctx *MetaDataHandlerContext) error {
	path, err := handler.ImageConverter.Convert(*ctx.media, imageconvertion.DEFAULT_BIG_MEDIAFORMAT)
	if err != nil {
		return err
	}

	ctx.media.Path = path
	return nil
}
