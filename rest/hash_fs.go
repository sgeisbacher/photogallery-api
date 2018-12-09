package rest

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/sgeisbacher/photogallery-api/media"
)

const FALLBACK_IMAGE_RELPATH = "fallback.jpg"

type HashFileSystem struct {
	DataRoot string
}

func (hfs HashFileSystem) Open(name string) (http.File, error) {
	groups := REGEXP_DATA_URL.FindStringSubmatch(name)
	if len(groups) != 3 {
		return nil, errors.New(fmt.Sprintf("url-path '%v' does not match", name))
	}

	mediaType := groups[1]
	hash := groups[2]

	media, err := media.Find(hash)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("could not find media '%v' in db", hash))
	}

	relPath := ""
	switch mediaType {
	case "thumb":
		relPath = media.ThumbnailPath
	case "big":
		relPath = media.Path
	case "orig":
		relPath = media.OrigPath
	default:
		return nil, errors.New(fmt.Sprintf("unknown mediaType '%v'", mediaType))
	}

	relPath = strings.TrimSpace(relPath)
	if len(relPath) == 0 {
		relPath = FALLBACK_IMAGE_RELPATH
	}

	// no need for security-checks on path, cause name is regexed and DataRoot is no user-input
	absPath := filepath.Join(hfs.DataRoot, relPath)
	file, err := os.Open(absPath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error while opening file '%v'", absPath))
	}
	fmt.Printf("returning file '%v'\n", file)
	return file, nil
}
