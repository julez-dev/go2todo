package fileutil

import "os"

func FileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func OpenAndTruncate(path string) (*os.File, error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)

	if err != nil {
		return nil, err
	}

	err = f.Truncate(0)

	if err != nil {
		return nil, err
	}

	return f, nil
}
