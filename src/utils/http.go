package utils

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"regexp"
)

func Request(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	return resp, nil
}

func DownloadFile(fileURL string, Version Version) (string, string, error) {
	resp, err := Request(fileURL)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if err := os.MkdirAll("out", os.ModePerm); err != nil {
		return "", "", fmt.Errorf("error creating output directory: %v", err)
	}

	filename := getFilename(resp)
	version := regexp.MustCompile(Version.Pattern).FindString(filename)
	if version == Version.Latest {
		return "", "", nil
	}

	println("Downloading", filename, "...")
	file, err := os.Create(filepath.Join("out", filename))
	if err != nil {
		return "", "", fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("error writing to file: %v", err)
	}

	fmt.Println("File downloaded successfully as", filename)
	return filename, version, nil
}

func getFilename(resp *http.Response) string {
	contentDisposition := resp.Header.Get("Content-Disposition")
	if contentDisposition != "" {
		parts := strings.Split(contentDisposition, ";")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if strings.HasPrefix(part, "filename=") {
				filename := strings.TrimPrefix(part, "filename=")
				return strings.Trim(filename, `"`)
			}
		}
	}
	return getDefaultFilename(resp)
}

func getDefaultFilename(resp *http.Response) string {
	u, err := url.Parse(resp.Request.URL.String())
	if err == nil {
		segments := strings.Split(u.Path, "/")
		return segments[len(segments)-1]
	}
	return "default_filename"
}
