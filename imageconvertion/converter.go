package imageconvertion

import "github.com/sgeisbacher/photogallery-api/media"

type ImageConverter interface {
	Convert(media media.Media, mediaFormat MediaFormat) (string, error)
}

type MediaFormat struct {
	Height int
	Width  int
	Path   string
}

var DEFAULT_THUMBNAIL_MEDIAFORMAT = MediaFormat{120, 120, "./data/thumbs/"}
var DEFAULT_BIG_MEDIAFORMAT = MediaFormat{800, 600, "./data/big/"}
