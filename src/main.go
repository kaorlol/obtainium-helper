package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"obtainium-helper/src/helpers"
	"obtainium-helper/src/utils"
)

var token string

func main() {
	token = getTokenArgs()

	fmt.Println("Started checking for updates")
	settings := utils.GetSettings()
	checkForUpdate(settings, settings.WaitTime)
	fmt.Println("Finished checking for updates")
}

func checkForUpdate(settings utils.Settings, waitTime int) {
	updateFound := false
	for name, download := range settings.ToDownload {
		fmt.Printf("Checking for updates for %s\n", name)
		var fileURL string
		var err error

		switch download.Type {
		case "web":
			fileURL, err = helpers.FetchURL(download.URL, download.Patterns, download.UrlEncoded, download.Agent)
		case "enumeration":
			ctx := context.Background()
			fileURL, err = helpers.Enumerate(ctx, download)
		case "artifact":
			if token == "" {
				panic("token is required for artifact downloads")
			}

			helpers.SetClient(token)
			fileURL, err = helpers.GetArtifacts(download.URL)
		}

		if err != nil {
			fmt.Printf("Error processing %s for %s: %v\n", download.Type, name, err)
			continue
		}

		if updateFound, settings = downloadApp(fileURL, name, download, settings); updateFound {
			return
		}
	}

	if !updateFound {
		fmt.Printf("No updates found, checking again in %d seconds\n", waitTime)
		time.Sleep(time.Duration(waitTime) * time.Second)
		checkForUpdate(settings, waitTime)
	}
}

func downloadApp(fileURL, name string, download utils.Download, settings utils.Settings) (bool, utils.Settings) {
	if fileURL == "" {
		return false, settings
	}

	appName, newIdentifier, err := utils.DownloadFile(fileURL, download)
	if err != nil {
		fmt.Printf("Error downloading file for %s: %v\n", name, err)
		return false, settings
	}

	if newIdentifier != "" {
		settings = utils.UpdateApp(settings, name, appName, newIdentifier)
		return true, settings
	}

	return false, settings
}

func getTokenArgs() string {
	tokenArgs := utils.Filter(os.Args[1:], func(arg string) bool {
		return strings.HasPrefix(arg, "--token=")
	})

	if len(tokenArgs) > 0 {
		return strings.TrimPrefix(tokenArgs[0], "--token=")
	}

	return ""
}
