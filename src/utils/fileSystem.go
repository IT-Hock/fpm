package utils

import (
	"archive/zip"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func GetExecutableName() string {
	return StripPath(os.Args[0])
}

func GetExecutablePath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}

	return filepath.Dir(ex), nil
}

func StripPath(path string) string {
	if path == "" {
		return ""
	}

	if strings.Contains(path, "/") {
		return path[strings.LastIndex(path, "/")+1:]
	}

	return path
}

func GetCacheDirectory() (string, error) {
	cacheDirectory := os.Getenv("XDG_CACHE_HOME")
	if cacheDirectory == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		cacheDirectory = path.Join(homeDir, ".cache")
	}

	return cacheDirectory, nil
}

func GetTempDirectory() (string, error) {
	tempDirectory := os.Getenv("TMPDIR")
	if tempDirectory == "" {
		tempDirectory = os.TempDir()
	}

	return tempDirectory, nil
}

func GetConfigDirectory() (string, error) {
	configDirectory := os.Getenv("XDG_CONFIG_HOME")
	if configDirectory == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configDirectory = path.Join(homeDir, ".config")
	}

	return configDirectory, nil
}

func GetFpmDirectory() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return path.Join(homeDir, ".local", "share", "fpm"), nil
}

func GetPackageDirectory() (string, error) {
	if os.Getenv("FPM_PACKAGE_DIR") != "" {
		return os.Getenv("FPM_PACKAGE_DIR"), nil
	}

	fpmDir, err := GetFpmDirectory()
	if err != nil {
		return "", err
	}

	return path.Join(fpmDir, "packages"), nil
}

func Unzip(source string, destination string) (string, error) {
	reader, err := zip.OpenReader(source)
	if err != nil {
		return "", err
	}
	defer func(reader *zip.ReadCloser) {
		err := reader.Close()
		if err != nil {
			panic(err)
		}
	}(reader)

	var rootDirectory string
	var firstDirectorySkipped = false

	for _, file := range reader.File {
		if strings.Contains(file.Name, "..") {
			continue
		}

		destinationFilePath := filepath.Join(destination, file.Name)
		destinationFilePath = strings.Replace(destinationFilePath, rootDirectory, "", -1)

		if file.FileInfo().IsDir() {
			if rootDirectory == "" {
				rootDirectory = file.Name
				if !firstDirectorySkipped {
					firstDirectorySkipped = true
					continue
				}
			}
			err := os.MkdirAll(destinationFilePath, file.Mode())
			if err != nil {
				return "", err
			}
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return "", err
		}

		targetFile, err := os.OpenFile(destinationFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return "", err
		}

		_, err = io.Copy(targetFile, fileReader)
		if err != nil {
			return "", err
		}

		err = targetFile.Close()
		if err != nil {
			return "", err
		}

		err = fileReader.Close()
		if err != nil {
			return "", err
		}
	}

	return rootDirectory, nil
}

// Exists returns whether a file or path exists
func Exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

// CreateDirectory creates a directory if it doesn't exist
func CreateDirectory(path string, mode os.FileMode) error {
	if Exists(path) {
		return nil
	}

	return os.MkdirAll(path, mode)
}

func MoveDirectory(source string, destination string) error {
	err := os.MkdirAll(destination, 0770)
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(source)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			err := MoveDirectory(source+"/"+entry.Name(), destination+"/"+entry.Name())
			if err != nil {
				return err
			}
		} else {
			err := Move(source+"/"+entry.Name(), destination+"/"+entry.Name())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Move moves a file from a source path to a destination path
// This must be used across the codebase for compatibility with Docker volumes
// and Golang (fixes Invalid cross-device link when using os.Rename)
func Move(sourcePath, destPath string) error {
	sourceAbs, err := filepath.Abs(sourcePath)
	if err != nil {
		return err
	}
	destAbs, err := filepath.Abs(destPath)
	if err != nil {
		return err
	}
	if sourceAbs == destAbs {
		return nil
	}
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}

	destDir := filepath.Dir(destPath)
	if !Exists(destDir) {
		err = os.MkdirAll(destDir, 0770)
		if err != nil {
			return err
		}
	}
	outputFile, err := os.Create(destPath)
	if err != nil {
		err := inputFile.Close()
		if err != nil {
			return err
		}
		return err
	}

	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		return err
	}

	err = inputFile.Close()
	if err != nil {
		return err
	}
	err = outputFile.Close()
	if err != nil {
		return err
	}
	if err != nil {
		if errRem := os.Remove(destPath); errRem != nil {
			return errRem
		}
		return err
	}

	err = os.Remove(sourcePath)
	if err != nil {
		return err
	}

	return nil
}

func GetFileName(path string) string {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path[strings.LastIndex(path, "/")+1:]
	}

	return filepath.Base(path)
}
