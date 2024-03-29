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

func downloadFileLessMemory(uri string) error {
	base := filepath.Base(uri)
	resp, err := http.Get(uri)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fileHandle, err := os.OpenFile(base, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer fileHandle.Close()
	_, err = io.Copy(fileHandle, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func HttpGet() (string, string, error) {
	var osUrl string
	var fileName string

	osName := runtime.GOOS

	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	var result entity.DownloadGithub
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Can not unmarshal JSON")
		return "", "", err
	}
	if osName == "linux" {
		osUrl = result.Assets[0].BrowserDownloadUrl
		fileName = result.Assets[0].Name
	}
	if osName == "darwin" {
		osUrl = result.Assets[1].BrowserDownloadUrl
		fileName = result.Assets[1].Name
	}
	if osName == "windows" {
		osUrl = result.Assets[2].BrowserDownloadUrl
		fileName = result.Assets[2].Name
	}
	err = downloadFileLessMemory(osUrl)
	if err != nil {
		return "", "", err
	}
	return fileName, result.TagName, nil
}

func LatestVersion() (string, error) {

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var result entity.DownloadGithub
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	return result.TagName, nil
}
