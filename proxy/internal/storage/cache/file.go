package cache

import "os"

// create a new file with the given path and content
func writeToFile(path string, content []byte) error {
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
func readFromFile(path string) ([]byte, error) {
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

// check if a file exists
func isFileExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return true
}

// count the number of files in a directory
func countFilesInDir(path string) (int, error) {
	// open the directory
	dir, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer dir.Close()

	// read the directory entries
	entries, err := dir.Readdir(-1)
	if err != nil {
		return 0, err
	}

	return len(entries), nil
}

// delete a file
func deleteFile(path string) error {
	// remove the file
	err := os.Remove(path)
	if err != nil {
		return err
	}

	return nil
}
