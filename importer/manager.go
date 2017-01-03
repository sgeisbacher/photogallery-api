package importer

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/sgeisbacher/photogallery-api/media"
)

type ImportManager struct {
	MediaService   *media.MediaService
	GalleryService *media.GalleryService
}

func (mgr ImportManager) ScanFolder(path string) {
	var wg sync.WaitGroup
	imageFilesChan := make(chan string)

	// start worker threads
	go mgr.handleImageFile(imageFilesChan, &wg)
	go mgr.handleImageFile(imageFilesChan, &wg)

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatalf("error while reading dir: %v", err)
	}
	for _, file := range files {
		wg.Add(1)
		if !file.IsDir() {
			imageFilesChan <- path + "/" + file.Name()
		}
	}

	close(imageFilesChan)
	wg.Wait()
}

func (mgr ImportManager) handleImageFile(imagesChan <-chan string, wg *sync.WaitGroup) {
	for filePath := range imagesChan {
		fileHash, err := hashFile(filePath)
		if err != nil {
			fmt.Printf("skipping file '%v' due to an error: %v\n", filePath, err)
			continue
		}
		fmt.Printf("md5ChkSum of '%v': %v\n", filePath, fileHash)
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
