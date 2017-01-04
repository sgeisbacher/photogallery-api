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
		bucket, err := tx.CreateBucketIfNotExists(BUCKET_GALLERIES)
		if err != nil {
			return err
		}

		gallery = getGalleryFromBucket(bucket, id)

		return err
	})

	return gallery, err
}

func (srv *GalleryService) Add(galleryName string) (*Gallery, error) {
	galleryId := BuildId(galleryName)
	gallery, _ := srv.FindGalleryById(galleryId)
	if gallery != nil {
		return gallery, nil
	}

	err := srv.Db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(BUCKET_GALLERIES)
		if err != nil {
			fmt.Println("error while creating/getting bucket:", err)
			return err
		}

		gallery = &Gallery{}
		gallery.Id = galleryId
		gallery.Name = galleryName

		galleryEncoded, err := gallery.gobEncode()
		return bucket.Put([]byte(gallery.Id), galleryEncoded)
	})

	return gallery, err
}

func (srv *GalleryService) AddMediaToGallery(galleryName string, media Media) error {
	srv.Add(galleryName)
	return srv.Db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(BUCKET_GALLERIES)
		if err != nil {
			fmt.Println("error while creating/getting bucket:", err)
			return err
		}

		galleryId := BuildId(galleryName)
		gallery := getGalleryFromBucket(bucket, galleryId)
		if gallery == nil {
			return errors.New(fmt.Sprintf("could not find gallery '%v'", galleryId))
		}

		gallery.Photos = append(gallery.Photos, media.Hash)

		galleryEncoded, err := gallery.gobEncode()
		return bucket.Put([]byte(gallery.Id), galleryEncoded)
	})
}

func BuildId(name string) string {
	// TODO normalize name and replace special-chars with _
	id := name
	return id
}

func getGalleryFromBucket(bucket *bolt.Bucket, id string) *Gallery {
	data := bucket.Get([]byte(id))
	if data == nil {
		return nil
	}
	gallery, err := gobDecodeGallery(data)
	if err != nil {
		return nil
	}
	return gallery
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
