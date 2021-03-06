package importer

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/sgeisbacher/photogallery-api/labels"
	"github.com/sgeisbacher/photogallery-api/media"
)

type ImportMediaData struct {
	name        string
	size        int64
	path        string
	galleryName string
}

type ImportManager struct {
}

func (mgr ImportManager) ScanFolder(path string) error {
	fmt.Printf("running import in '%v' ...\n", path)
	var wg sync.WaitGroup
	imagesChan := make(chan ImportMediaData)

	// start worker threads
	go mgr.handleImageFile(imagesChan, &wg)
	go mgr.handleImageFile(imagesChan, &wg)

	path = addSlash(path)

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return errors.New(fmt.Sprintf("error while reading dir: %v", err))
	}
	for _, file := range files {
		if !file.IsDir() {
			fmt.Printf("skipping file '%v' (not allowed here), because its on gallery-folder-level\n", path+file.Name())
			continue
		}
		if file.Name() == "data" {
			fmt.Printf("skipping file '%v' (not allowed here), because its on ignore-list\n", path+file.Name())
			continue
		}
		fmt.Printf("importing folder %q ...\n", file.Name())
		scanGalleryFolder(file.Name(), addSlash(path+file.Name()), imagesChan, &wg)
	}

	close(imagesChan)
	wg.Wait()

	fmt.Println("Importer ... done")
	return nil
}

func scanGalleryFolder(galleryName, path string, imagesChan chan ImportMediaData, wg *sync.WaitGroup) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Printf("error while reading dir: %v\n", err)
		return
	}
	for _, file := range files {
		if file.IsDir() {
			fmt.Printf("skipping directory '%v' (not allowed here), because its on image-level\n", addSlash(path+file.Name()))
			continue
		}
		wg.Add(1)
		importMediaData := ImportMediaData{
			name:        file.Name(),
			path:        path + file.Name(),
			galleryName: galleryName,
			size:        file.Size(),
		}
		imagesChan <- importMediaData
	}
}

func (mgr ImportManager) handleImageFile(imagesChan <-chan ImportMediaData, wg *sync.WaitGroup) {
	for importMediaData := range imagesChan {
		fileHash, err := hashFile(importMediaData.path)
		if err != nil {
			fmt.Printf("skipping file '%v' due to an error: %v\n", importMediaData.path, err)
			continue
		}

		m := media.Media{
			Hash:      fileHash,
			Name:      importMediaData.name,
			Path:      fmt.Sprintf("%s/%s", importMediaData.galleryName, importMediaData.name),
			OrigPath:  importMediaData.path,
			Size:      importMediaData.size,
			MediaType: media.MEDIA_TYPE_PHOTO,
		}
		label := labels.Add(importMediaData.galleryName)
		media.Add(&m)
		err = labels.LabelMedia(label.ID, m)
		if err != nil {
			fmt.Printf("error while putting media to gallery '%v': %v\n", importMediaData.galleryName, err)
			// TODO delete media!!!
		}
		wg.Done()
	}
}

func hashFile(filename string) (string, error) {
	var md5ChkSum string
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("could not open file:", filename)
		return md5ChkSum, err
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		fmt.Println("could not hash file:", filename)
		return md5ChkSum, err
	}
	md5ChkSum = hex.EncodeToString(hash.Sum(nil))
	return md5ChkSum, nil
}

func addSlash(path string) string {
	path = strings.TrimSpace(path)
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	return path
}
