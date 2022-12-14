package main

import (
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
)

// get directory size, in B
func DirSize(path string) (float64, error) {
	var size float64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += float64(info.Size())
		}
		return err
	})
	return size, err
}

// remove a directory
func RemoveDir(dir []string) error {
	log.Println(dir)
	for _, d := range dir {
		err := os.RemoveAll(d)
		if err != nil {
			log.Println("Unable to remove dir: ", d, " got error: ", err)
		}
	}
	return nil
}

// get a map with directories and their sizes, in GB
var (
	kb = math.Pow(10, 3)
	mb = math.Pow(10, 6)
	gb = math.Pow(10, 9)
	tb = math.Pow(10, 12)
)

func getDirNames(histDir string) map[string]float64 {
	dir := histDir
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalf("Got error: %v", err)
	}
	dirNames := map[string]float64{}
	for _, file := range files {
		sub := dir + "/" + file.Name()
		size, err := DirSize(sub)
		dirNames[sub] = size / gb
		if err != nil {
			log.Fatalf("Got error when scanning size of directories: %e", err)
		}
		log.Println(sub, size/gb)
	}
	return dirNames
}
