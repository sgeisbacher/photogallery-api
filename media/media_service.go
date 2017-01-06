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
	Db             *bolt.DB
	GalleryService *GalleryService
}

func (srv *MediaService) FindMediaByHash(hash string) (*Media, error) {
	var media *Media
	err := srv.Db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(BUCKET_MEDIAS)
		if bucket == nil {
			return errors.New(fmt.Sprintf("bucket '%v' not found", BUCKET_MEDIAS))
		}
		data := bucket.Get([]byte(hash))
		if data == nil {
			return nil
		}
		var err error
		media, err = gobDecodeMedia(data)
		return err
	})

	return media, err
}

func (srv *MediaService) Add(media Media) error {
	if media, _ := srv.FindMediaByHash(media.Hash); media != nil {
		fmt.Printf("skipping '%v' (%v), it already exists\n", media.Path, media.Hash)
		return nil
	}
	return srv.Db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(BUCKET_MEDIAS)
		if err != nil {
			fmt.Println("error while creating/getting bucket 'medias':", err)
			return err
		}

		mediaEncoded, err := media.gobEncode()
		err = bucket.Put([]byte(media.Hash), mediaEncoded)

		fmt.Printf("added media '%v' with md5ChkSum: %v\n", media.Path, media.Hash)

		return err
	})
}

func (srv *MediaService) FindAll() ([]*Media, error) {
	var medias []*Media
	err := srv.Db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(BUCKET_MEDIAS)
		if bucket == nil {
			return errors.New(fmt.Sprintf("bucket '%v' not found", string(BUCKET_MEDIAS)))
		}
		cursor := bucket.Cursor()
		for key, data := cursor.First(); data != nil; key, data = cursor.Next() {
			media, err := gobDecodeMedia(data)
			if err != nil {
				fmt.Printf("could not decode '%v': %v\n", key, err)
				continue
			}
			medias = append(medias, media)
		}
		return nil
	})
	return medias, err
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
