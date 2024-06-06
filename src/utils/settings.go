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
	Name       string     `json:"name"`
	Type       string     `json:"type"`
	Identifier Identifier `json:"identifier"`
	URL        string     `json:"url"`
	Patterns   []string   `json:"patterns"`
	UrlEncoded bool       `json:"url_encoded"`
}

type Identifier struct {
	Latest         string `json:"latest"`
	EnumLimit      int    `json:"enum_limit"`
	IncrementLimit int    `json:"increment_limit"`
	Pattern        string `json:"pattern"`
}

func GetSettings() Settings {
	file := filepath.Join("settings.json")
	data, _ := os.ReadFile(file)

	var settings Settings
	json.Unmarshal(data, &settings)
	return settings
}

func UpdateApp(settings Settings, name, appName, identifier string) Settings {
	download := settings.ToDownload[name]
	download.Name = appName
	download.Identifier.Latest = identifier
	settings.ToDownload[name] = download

	file := filepath.Join("settings.json")
	data, _ := json.Marshal(settings)
	os.WriteFile(file, data, 0o644)
	return settings
}
