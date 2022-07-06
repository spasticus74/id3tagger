package fileops

import (
	"errors"
	"os"
	"path/filepath"
)

func CheckPath(dirPath string) error {
	// first, check that what we've been given actually exists
	theDir, err := os.Open(dirPath)

	if err != nil {
		theDir.Close()
		return err
	} else {
		fileInfo, err := theDir.Stat()
		if err != nil {
			theDir.Close()
			return err
		}

		// check that the path given points to a directory and not a file or some other thing
		if !fileInfo.IsDir() {
			theDir.Close()
			return errors.New("the given path exists but is not a directory")
		}
	}

	// all clear
	theDir.Close()
	return nil
}

func GetMP3sInDir(dirPath string) ([]string, error) {
	// open the dir
	theDir, _ := os.Open(dirPath)
	defer theDir.Close()

	// read the dir
	files, err := theDir.Readdir(0)
	if err != nil {
		//fmt.Println(err)
		return nil, err
	}

	// gather the ones that are mp3
	var ret []string
	for _, v := range files {
		if filepath.Ext(v.Name()) == ".mp3" {
			ret = append(ret, v.Name())
		}
	}

	return ret, nil
}
