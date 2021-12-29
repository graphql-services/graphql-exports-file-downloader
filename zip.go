package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"io"
	"net/http"
)

// ZipFiles compresses one or many files into a single zip archive file.
// Param 1: filename is the output zip file's name.
// Param 2: files is a list of files to add to the zip.
func ZipFiles(files []FileMetadataFile, progressFn func(progress float32) error) (buf bytes.Buffer, err error) {

	// newZipFile, err := os.Create(filename)
	// if err != nil {
	// 	return err
	// }
	// defer newZipFile.Close()
	// var buf bytes.Buffer
	zipBuffer := bufio.NewWriter(&buf)
	zipWriter := zip.NewWriter(zipBuffer)
	defer zipWriter.Close()

	// Add files to zip
	for i, file := range files {
		err = AddFileToZip(zipWriter, file)
		if err != nil {
			return
		}
		err = progressFn(float32(i+1) / float32(len(files)))
		if err != nil {
			return
		}
	}
	return buf, nil
}

func AddFileToZip(zipWriter *zip.Writer, file FileMetadataFile) (err error) {
	response, err := http.Get(file.URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	write, err := zipWriter.CreateHeader(&zip.FileHeader{
		Name:   file.Name,
		Method: zip.Deflate,
	})
	if err != nil {
		return
	}

	io.Copy(write, response.Body)

	// fileToZip, err := os.Open(filename)
	// if err != nil {
	// 	return err
	// }
	// defer fileToZip.Close()

	// // Get the file information
	// info, err := fileToZip.Stat()
	// if err != nil {
	// 	return err
	// }

	// header, err := zip.FileInfoHeader(info)
	// if err != nil {
	// 	return err
	// }

	// // Using FileInfoHeader() above only uses the basename of the file. If we want
	// // to preserve the folder structure we can overwrite this with the full path.
	// header.Name = filename

	// // Change to deflate to gain better compression
	// // see http://golang.org/pkg/archive/zip/#pkg-constants
	// header.Method = zip.Deflate

	// writer, err := zipWriter.CreateHeader(header)
	// if err != nil {
	// 	return err
	// }
	// _, err = io.Copy(writer, fileToZip)
	// return err
	return
}
