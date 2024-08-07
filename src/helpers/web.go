package helpers

import (
	"fmt"
	"io"
	"net/url"

	"obtainium-helper/src/utils"
)

func fetchText(url string, agent *string) (string, error) {
	resp, err := utils.Request(url, agent)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func FetchURL(URL string, patterns []string, app utils.Download) (string, error) {
	if len(patterns) == 0 {
		return URL, nil
	}

	text, err := fetchText(URL, app.Agent)
	if err != nil {
		return "", err
	}

	if app.UrlEncoded {
		text, err = url.QueryUnescape(text)
		if err != nil {
			return "", err
		}
	}

	matches := utils.FindAllPatternsInText(text, patterns[0])
	if len(matches) == 0 {
		return "", fmt.Errorf("pattern not found in text: %s", patterns[0])
	}

	if len(patterns) == 1 {
		highestVersion := utils.ExtractIdentifier(matches[0], app.Identifier.Pattern)
		if highestVersion == "" {
			return "", fmt.Errorf("identifier not found in match: %s", matches[0])
		}
		highestMatch := matches[0]

		for _, match := range matches[1:] {
			version := utils.ExtractIdentifier(match, app.Identifier.Pattern)
			if version == "" {
				continue
			}

			compare, err := utils.CompareVersions(version, highestVersion)
			if err != nil {
				return "", err
			}

			if compare > 0 {
				highestVersion = version
				highestMatch = match
			}
		}

		return highestMatch, nil
	}

	return FetchURL(matches[0], patterns[1:], app)
}
