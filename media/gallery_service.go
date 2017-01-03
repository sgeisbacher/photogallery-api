package media

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"

	"github.com/boltdb/bolt"
)

var (
	BUCKET_GALLERIES = []byte("galleries")
)

type Gallery struct {
	Id     string
	Name   string
	Year   int
	Photos []string
}

type GalleryService struct {
	Db *bolt.DB
}

func (srv *GalleryService) FindGalleryById(id string) (*Gallery, error) {
	var gallery *Gallery
	err := srv.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(BUCKET_GALLERIES)
		data := b.Get([]byte(id))
		if data == nil {
			return errors.New(fmt.Sprintf("gallery '%v' not found!", id))
		}
		var err error
		gallery, err = gobDecodeGallery(data)
		return err
	})

	return gallery, err
}

func (srv *GalleryService) Add(gallery Gallery) error {
	return srv.Db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(BUCKET_GALLERIES)
		if err != nil {
			return err
		}

		if len(gallery.Id) == 0 {
			gallery.BuildId()
		}
		galleryEncoded, err := gallery.gobEncode()
		return bucket.Put([]byte(gallery.Id), galleryEncoded)
	})
}

func (gallery *Gallery) BuildId() string {
	// TODO normalize name and replace special-chars with _
	id := gallery.Name
	gallery.Id = id
	return id
}

func (gallery Gallery) gobEncode() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(gallery)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func gobDecodeGallery(data []byte) (*Gallery, error) {
	var g *Gallery
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&g)
	if err != nil {
		return nil, err
	}
	return g, nil
}
