package importer

import (
	"fmt"
	"io/ioutil"
	"log"
	"sync"
)

type ImportManager struct {
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
		imageFilesChan <- file.Name()
	}

	close(imageFilesChan)
	wg.Wait()
}

func (mgr ImportManager) handleImageFile(imagesChan <-chan string, wg *sync.WaitGroup) {
	for filename := range imagesChan {
		fmt.Println("processing", filename, "...")
		wg.Done()
	}
}
