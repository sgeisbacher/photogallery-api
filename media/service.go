package media

import (
	"fmt"
	"time"

	cache "github.com/sgeisbacher/go-cache"
)

const (
	MEDIA_TYPE_PHOTO = iota
	MEDIA_TYPE_VIDEO = iota
)

type Media struct {
	Hash          string
	Name          string
	Size          int64
	ThumbnailPath string
	Path          string
	OrigPath      string
	MediaType     int
	ShootTime     time.Time
}

var mediaIdx *cache.Cache = cache.New(cache.NoExpiration, 0)
var mediaLabelsIdx *cache.Cache = cache.New(cache.NoExpiration, 0)

// Find finds
func Find(hash string) (*Media, error) {
	media, found := mediaIdx.Get(hash)
	if !found {
		return nil, fmt.Errorf("could not find media %q", hash)
	}
	return media.(*Media), nil
}

// Add adds
func Add(media *Media) {
	mediaIdx.Set(media.Hash, media, cache.NoExpiration)
	mediaLabelsIdx.Set(media.Hash, []string{}, cache.NoExpiration)
}

// FindAll finds all
func FindAll() []*Media {
	var medias []*Media
	for _, media := range mediaIdx.Items() {
		medias = append(medias, media.Object.(*Media))
	}
	return medias
}

func GetLabels(media *Media) ([]string, error) {
	labels, found := mediaLabelsIdx.Get(media.Hash)
	if !found {
		return nil, fmt.Errorf("could not find labels for media %q, media not found", media.Hash)
	}
	return labels.([]string), nil
}

func AddLabel(media *Media, labelID string) error {
	err := mediaLabelsIdx.Append(media.Hash, labelID)
	if err != nil {
		return fmt.Errorf("could not add label %q to media %q: %v", labelID, media.Hash, err)
	}
	return nil
}
