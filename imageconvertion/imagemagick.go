package imageconvertion

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/sgeisbacher/goutils/fileutils"
	"github.com/sgeisbacher/photogallery-api/media"
)

type ImageMagickImageConverter struct{}

func (cvtr ImageMagickImageConverter) Convert(media media.Media, mediaFormat MediaFormat) (string, error) {
	path, err := exec.LookPath("convert")
	if err != nil {
		return "", err
	}

	thumbnailPath := fileutils.AppendSlash(mediaFormat.Path)

	err = os.MkdirAll(thumbnailPath, 0777)
	if err != nil {
		return "", err
	}

	thumbnailPath = thumbnailPath + media.Hash + ".jpg"

	dimension, dErr := BuildDimension(mediaFormat)
	if dErr != nil {
		return "", dErr
	}

	convertCmd := exec.Command(path, "-resize", dimension, media.OrigPath, thumbnailPath)
	convertCmd.Stdout = os.Stdout
	convertCmd.Stderr = os.Stderr
	err = convertCmd.Run()
	if err != nil {
		return "", err
	}

	return thumbnailPath, nil
}

func BuildDimension(mediaFormat MediaFormat) (string, error) {
	if mediaFormat.Height == 0 && mediaFormat.Width == 0 {
		return "", errors.New("invalid mediaformat, both width and height not set, only one allow")
	}

	heightStr := ""
	if mediaFormat.Height > 0 {
		heightStr = strconv.Itoa(mediaFormat.Height)
	}

	widthStr := ""
	if mediaFormat.Width > 0 {
		widthStr = strconv.Itoa(mediaFormat.Width)
	}

	return fmt.Sprintf("%vx%v", heightStr, widthStr), nil
}
