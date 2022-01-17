package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// learningContainerDownload will create a downloads folder relative to the
// applications working directory. It then later calls the downloadFile func and
// passes in the filepath and url specified
func learningContainerDownload() {

	//MkdirAll will not error if directory already exists
	err1 := os.MkdirAll("downloads", 0744)
    if err1 != nil {
        panic(err1)
	}

	fileUrl := "https://www.learningcontainer.com/wp-content/uploads/2020/04/sample-text-file.txt"
	err2 := downloadFile("downloads/sample-text-file.txt", fileUrl)
	if err2 != nil {
		panic(err2)
	}
	fmt.Println("Downloaded: " + fileUrl)
}

// downloadFile will take the values passed in and perform a Get and Create of the file
func downloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err!= nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}