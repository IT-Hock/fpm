package utils

import (
	"io"
	"net/http"
	"os"
)

func DownloadFile(url string, path string) error {

	// Create the file
	out, err := os.Create(path)
	if err != nil {
		return err
	}

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			panic(err)
		}
	}(out)

	defer func(resp *http.Response) {
		err := resp.Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp)

	return nil
}
