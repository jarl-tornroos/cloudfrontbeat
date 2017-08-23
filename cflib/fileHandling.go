package cflib

import (
	"compress/gzip"
	"io/ioutil"
	"os"
)

// FileHandling makes it easy to read from and write to files
type FileHandling struct {
	File string
}

// SetFile, setter for file
func (f *FileHandling) SetFile(file string) *FileHandling {
	f.File = file
	return f
}

// ReadContent from a normal text file
func (f *FileHandling) ReadContent() (string, error) {
	dat, err := ioutil.ReadFile(f.File)
	if err != nil {
		return "", err
	}
	content := string(dat)

	return content, nil
}

// ReadZipContent from a compressed file
func (f *FileHandling) ReadZipContent() (string, error) {
	file, err := os.Open(f.File)
    if err != nil {
        return "", err
    }
    defer file.Close()

    fileUncompressed, err := gzip.NewReader(file)
    if err != nil {
        return "", err
    }
    defer fileUncompressed.Close()

    dat, err := ioutil.ReadAll(fileUncompressed)
    if err != nil {
        return "", err
    }

	content := string(dat)

    return content, nil
}

// Write append content to a file
func (f *FileHandling) Write(add string) (int, error) {

	var file *os.File
	var err error

	// Create and open file if it doesn't exist, otherwise open it for appending more content
	if _, err = os.Stat(f.File); os.IsNotExist(err) {
		file, err = os.Create(f.File)
	} else {
		file, err = os.OpenFile(f.File, os.O_APPEND|os.O_WRONLY, 0644)
	}

	// Close file after the function return
	defer file.Close()

	bytes, err := file.WriteString(add + "\n")

	file.Sync()

	return bytes, err
}

// Delete a file
func (f *FileHandling) Delete() (bool) {
	err := os.Remove(f.File)
	if err != nil {
		return false
	}
	return true
}
