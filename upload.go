package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

type FileUpload struct {
	Filename    string
	Size        int
	ContentType string
}

type GetUploadURLResult struct {
	ID        string
	UploadURL string
}

// ZipFiles compresses one or many files into a single zip archive file.
// Param 1: filename is the output zip file's name.
// Param 2: files is a list of files to add to the zip.
func Upload(filename string, buf bytes.Buffer) (fileId string, err error) {
	response, err := getUploadURL(FileUpload{
		Filename:    filename,
		Size:        len(buf.Bytes()),
		ContentType: "application/zip",
	})
	if err != nil {
		return
	}
	fileId = response.ID

	req, err := http.NewRequest("PUT", response.UploadURL, bytes.NewReader(buf.Bytes()))
	if err != nil {
		return
	}

	_, err = http.DefaultClient.Do(req)
	return
}

func getUploadURL(fu FileUpload) (result GetUploadURLResult, err error) {
	uploadURL := os.Getenv("FILES_UPLOAD_URL")
	if uploadURL == "" {
		panic("Missing 'FILES_UPLOAD_URL' environment variable")
	}

	reqBody, err := json.Marshal(fu)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", uploadURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return
	}
	req.Header.Set("content-type", "application/json")
	client := getClientcredentialsClient()
	res, err := client.Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &result)
	return
}
