package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"ravro_dcrpt/entity"
	"runtime"
)

const url string = "https://api.github.com/repos/ravro-ir/ravro_dcrpt/releases/latest"

func downloadFileLessMemory(uri string) {
	base := filepath.Base(uri)
	resp, err := http.Get(uri)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fileHandle, err := os.OpenFile(base, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer fileHandle.Close()
	_, err = io.Copy(fileHandle, resp.Body)
	if err != nil {
		panic(err)
	}
}

func HttpGet() {
	// Get request
	var osUrl string
	osName := runtime.GOOS

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("No response from request")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body) // response body is []byte

	var result entity.DownloadGithub
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to the go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}
	if osName == "linux" {
		osUrl = result.Assets[0].BrowserDownloadUrl
	}
	if osName == "darwin" {
		osUrl = result.Assets[1].BrowserDownloadUrl
	}
	if osName == "windows" {
		osUrl = result.Assets[2].BrowserDownloadUrl
	}
	downloadFileLessMemory(osUrl)
	fmt.Println("[++++] The latest version file downloaded ")
}
