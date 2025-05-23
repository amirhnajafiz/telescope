package cache

import (
	"os"
	"path/filepath"
)

// write create a new file with the given path and content
func write(path string, content []byte) error {
	// create directories in the path if they do not exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	// open the file for writing
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// write the content to the file
	_, err = file.Write(content)
	if err != nil {
		return err
	}

	return nil
}

// read the content of a file
func read(path string) ([]byte, error) {
	// open the file for reading
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// read the content from the file
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return content, nil
}

// exists, check if a file exists
func exists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return true
}
