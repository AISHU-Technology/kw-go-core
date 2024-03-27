package filex

import (
	"errors"
	"github.com/AISHU-Technology/kw-go-core/idx"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

const (
	FilePoint     string = "."
	FileCharacter string = "/"
)

func GetNewFileName(fileSuffix string) string {
	return idx.NewUUID().String() + FilePoint + fileSuffix
}

func GetFileName(path string) string {
	return GetFileNameSuffix(path)[:GetFileLastIndex(path, FilePoint)]
}

func GetFileSuffix(path string) string {
	return path[GetFileLastIndex(path, FilePoint)+1:]
}

func GetFileNameSuffix(path string) string {
	return path[:GetFileLastIndex(path, FileCharacter)]
}

func GetFileLastIndex(path, character string) int {
	index := strings.LastIndex(path, character)
	return index
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return false
}

func CreateFile(path string) bool {
	file, err := os.Create(path)
	if err != nil {
		return false
	}
	defer file.Close()
	return true
}

func CreateDir(absPath string) error {
	return os.MkdirAll(absPath, os.ModePerm)
}

func IsDir(path string) bool {
	file, err := os.Stat(path)
	if err != nil {
		return false
	}
	return file.IsDir()
}

func RemoveFile(path string) error {
	return os.Remove(path)
}

func CopyFile(src string, desc string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	distFile, err := os.Create(desc)
	if err != nil {
		return err
	}
	defer distFile.Close()

	var tmp = make([]byte, 1024*4)
	for {
		n, err := srcFile.Read(tmp)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		_, err = distFile.Write(tmp[:n])
		if err != nil {
			return err
		}
	}
}

// ListDir lists all the file or dir names in the specified directory.
// Note that ListDir don't traverse recursively.
func ListDir(dirname string) ([]string, error) {
	infos, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(infos))
	for i, info := range infos {
		names[i] = info.Name()
	}
	return names, nil
}

// GetDirAllFilePaths gets all the file paths in the specified directory recursively.
func GetDirAllFilePaths(dirname string) ([]string, error) {
	// Remove the trailing path separator if dirname has.
	dirname = strings.TrimSuffix(dirname, FileCharacter)
	infos, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	paths := make([]string, 0, len(infos))
	for _, info := range infos {
		path := dirname + FileCharacter + info.Name()
		if info.IsDir() {
			tmp, err := GetDirAllFilePaths(path)
			if err != nil {
				return nil, err
			}
			paths = append(paths, tmp...)
			continue
		}
		paths = append(paths, path)
	}
	return paths, nil
}
