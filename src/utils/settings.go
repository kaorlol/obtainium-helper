package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Settings struct {
	ToDownload map[string]Download `json:"to_download"`
	WaitTime   int                 `json:"wait_time"`
}

type Download struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	Version    Version  `json:"version"`
	URL        string   `json:"url"`
	Patterns   []string `json:"patterns"`
	UrlEncoded bool     `json:"url_encoded"`
}

type Version struct {
	Latest  string `json:"latest"`
	Pattern string `json:"pattern"`
}

func GetSettings() Settings {
	file := filepath.Join("settings.json")
	data, _ := os.ReadFile(file)

	var settings Settings
	json.Unmarshal(data, &settings)
	return settings
}

func UpdateApp(settings Settings, name, appName, version string) Settings {
	download := settings.ToDownload[name]
	download.Name = appName
	download.Version.Latest = version
	settings.ToDownload[name] = download

	file := filepath.Join("settings.json")
	data, _ := json.Marshal(settings)
	os.WriteFile(file, data, 0o644)
	return settings
}
