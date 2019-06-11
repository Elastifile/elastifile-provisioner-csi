package helputils

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-errors/errors"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}

// open specified path, and call fail() if failed.
func OpenFile(path string, fail func(err error)) *os.File {
	file, err := os.Open(path)
	if err != nil {
		fail(err)
		return nil
	}
	return file
}

// close specified fail, and panic if failed.
func MustCloseFile(file *os.File) {
	if err := file.Close(); err != nil {
		panic(err)
	}
}

func ReadAll(filepath string) ([]byte, error) {
	return ReadFrom(filepath, 0)
}

func ReadFrom(filepath string, offset int64) ([]byte, error) {
	file, err := os.OpenFile(filepath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer MustCloseFile(file)

	if _, err := file.Seek(offset, 0); err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func FindFilesByPattern(root string, pattern string) (files []string, err error) {
	filepath.Walk(root, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			if r, err := regexp.MatchString(pattern, path); err == nil && r {
				if rel, err := filepath.Rel(root, path); err != nil {
					return errors.New("Error extracting relative filepath: " + err.Error())
				} else {
					files = append(files, rel)
				}
			}
		}
		return nil
	})
	return files, nil
}

func GetWd() (string, error) {
	if cwd, err := os.Getwd(); err != nil {
		return "", errors.New(err)
	} else {
		return cwd, nil
	}
}

// WdRel returns the relative path to the current working directory.
// It utilizes os.Getwd() and filepath.Rel().
func WdRel(path string) (string, error) {
	workingDir, err := GetWd()
	if err != nil {
		return "", err
	}
	if rel, err := filepath.Rel(workingDir, path); err != nil {
		return "", errors.New(err)
	} else {
		return rel, nil
	}
}

func WriteJson(file string, config interface{}) error {
	jsonBody, err := json.Marshal(config)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(file, jsonBody, os.ModePerm)
	return err
}

func WriteFile(file string, data []byte, mode os.FileMode) error {
	if dir := filepath.Dir(file); dir != "" {
		err := os.MkdirAll(dir, mode)
		if err != nil {
			return err
		}
	}
	err := ioutil.WriteFile(file, data, mode)
	if err != nil {
		return err
	}
	err = ForceFileMode(file, mode)
	return err
}

func FullPath(root string, path string) string {
	if strings.HasPrefix(path, "/") {
		return path
	}
	return filepath.Join(root, path)
}
