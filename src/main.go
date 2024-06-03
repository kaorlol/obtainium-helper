package main

import (
	"fmt"
	"time"

	"obtainium-helper/src/helpers"
	"obtainium-helper/src/utils"
)

func main() {
	fmt.Println("Started checking for updates")
	settings := utils.GetSettings()
	checkForUpdate(settings, settings.WaitTime)

	fmt.Println("Finished checking for updates")
}

func checkForUpdate(settings utils.Settings, waitTime int) {
	updateFound := false

	for name, download := range settings.ToDownload {
		if download.Type == "web" {
			fileURL, err := helpers.FetchURL(download.URL, download.Patterns, download.UrlEncoded)
			if err != nil {
				fmt.Printf("Error fetching URL for %s: %v\n", name, err)
				continue
			}

			appName, newVersion, err := utils.DownloadFile(fileURL, download.Version)
			if err != nil {
				fmt.Printf("Error downloading file for %s: %v\n", name, err)
				continue
			}

			if newVersion != "" {
				settings = utils.UpdateApp(settings, name, appName, newVersion)
				updateFound = true
				break
			}
		}
	}

	if !updateFound {
		fmt.Printf("No updates found, checking again in %d seconds\n", waitTime)
		time.Sleep(time.Duration(waitTime) * time.Second)
		checkForUpdate(settings, waitTime)
	}
}
