package labels

import (
	"fmt"

	cache "github.com/sgeisbacher/go-cache"
	"github.com/sgeisbacher/goutils/datautils"
	"github.com/sgeisbacher/photogallery-api/media"
)

// Label it can group up medias
type Label struct {
	ID   string
	Name string
}

var labelsIdx *cache.Cache = cache.New(cache.NoExpiration, 0)
var labelMediasIdx *cache.Cache = cache.New(cache.NoExpiration, 0)

// Find searches for label with id
func Find(id string) (*Label, error) {
	label, ok := labelsIdx.Get(id)
	if !ok {
		return nil, fmt.Errorf("label %q not found", id)
	}
	return label.(*Label), nil
}

// GetMedias returns all media-ids labelled with labelID
func GetMedias(labelID string) ([]string, error) {
	medias, found := labelMediasIdx.Get(labelID)
	if !found {
		return nil, fmt.Errorf("could not find medias for label %q, label not found", labelID)
	}
	return medias.([]string), nil
}

// Add creates a new label
func Add(name string) *Label {
	id := datautils.ToID(name)
	label := &Label{id, name}
	err := labelsIdx.Add(id, label, cache.NoExpiration)
	if err == nil {
		labelMediasIdx.Set(id, []string{}, cache.NoExpiration)
	}
	return label
}

// LabelMedia labels the media with the given label
func LabelMedia(labelID string, m media.Media) error {
	err := labelMediasIdx.Append(labelID, m.Hash)
	if err != nil {
		return fmt.Errorf("could not label media %q with %q: %v", m.Hash, labelID, err)
	}
	return media.AddLabel(&m, labelID)
}

// FindAll gives all label objects without (!) their medias
func FindAll() []*Label {
	var labels []*Label
	for _, item := range labelsIdx.Items() {
		label := item.Object.(*Label)
		labels = append(labels, label)
	}
	return labels
}
