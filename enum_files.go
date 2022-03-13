package main

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	GET_FULL_PATHNAME = true
	GET_FILENAME_ONLY = false
)

// need ".ext" form
func GetFilesList(dirname string, ext string, getFullPathname bool) ([]string, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	names, err := f.Readdirnames(-1)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	Ext := strings.ToUpper(ext)

	var namesFiles []string
	for _, name := range names {
		filename := filepath.Join(dirname, name)
		fileinfo, err := os.Lstat(filename)
		if err == nil && !fileinfo.IsDir() {
			if strings.ToUpper(filepath.Ext(filename)) == Ext {
				if getFullPathname {
					namesFiles = append(namesFiles, filename) // append full path to file
				} else {
					namesFiles = append(namesFiles, name)
				}
			}
		}

	}
	sort.Strings(namesFiles)
	return namesFiles, nil
}

func ListFiles(dir string, exts []string, getFullPathname bool) []string {
	var files []string
	for _, ext := range exts {
		if f, err := GetFilesList(dir, ext, getFullPathname); err != nil {
			log.Fatalln(err)
		} else {
			files = append(files, f...)
		}
	}
	return files
}

// Split full file path+name to (dir, base, ext)
func SplitPath(path string) (string, string, string) {
	dir, file := filepath.Split(path)
	ext := filepath.Ext(path) // include dot (.mp3)
	base := strings.TrimSuffix(file, ext)
	return dir, base, ext
}

func FileName(path string) string {
	_, filename := filepath.Split(path)
	return filename
}
