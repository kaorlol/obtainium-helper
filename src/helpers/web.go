package helpers

import (
	"fmt"
	"io"
	"net/url"
	"regexp"

	"obtainium-helper/src/utils"
)

func fetchText(url string, agent *string) (string, error) {
	resp, err := utils.Request(url, agent)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func findPatternInText(text string, pattern string) (string, bool) {
	re := regexp.MustCompile(pattern)
	match := re.FindString(text)
	return match, match != ""
}

func FetchURL(URL string, patterns []string, urlEncoded bool, agent *string) (string, error) {
	if len(patterns) == 0 {
		return URL, nil
	}

	text, err := fetchText(URL, agent)
	if err != nil {
		return "", err
	}

	if urlEncoded {
		text, err = url.QueryUnescape(text)
		if err != nil {
			return "", err
		}
	}

	match, found := findPatternInText(text, patterns[0])
	if !found {
		return "", fmt.Errorf("pattern not found in Text: %s", patterns[0])
	}

	if len(patterns) == 1 {
		return match, nil
	}

	return FetchURL(match, patterns[1:], urlEncoded, agent)
}
