package media

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"

	"github.com/boltdb/bolt"
)

const (
	MEDIA_TYPE_PHOTO = iota
	MEDIA_TYPE_VIDEO = iota
)

var (
	BUCKET_MEDIAS = []byte("medias")
)

type Media struct {
	Hash          string
	Name          string
	Size          int
	ThumbnailPath string
	Path          string
	OrigPath      string
	MediaType     int
}

type MediaService struct {
	Db *bolt.DB
}

func (srv *MediaService) FindMediaByHash(hash string) (*Media, error) {
	var media *Media
	err := srv.Db.View(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(BUCKET_MEDIAS)
		if err != nil {
			return err
		}
		data := bucket.Get([]byte(hash))
		if data == nil {
			return errors.New(fmt.Sprintf("media '%v' not found!", hash))
		}
		media, err = gobDecodeMedia(data)
		return err
	})

	return media, err
}

func (srv *MediaService) Add(media Media) error {
	if media, _ := srv.FindMediaByHash(media.Hash); media != nil {
		fmt.Printf("skipping '%v', it already exists\n", media.Path)
		return nil
	}
	return srv.Db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(BUCKET_MEDIAS)
		if err != nil {
			return err
		}

		mediaEncoded, err := media.gobEncode()
		return bucket.Put([]byte(media.Hash), mediaEncoded)
	})
}

func (media Media) gobEncode() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(media)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func gobDecodeMedia(data []byte) (*Media, error) {
	var m *Media
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
