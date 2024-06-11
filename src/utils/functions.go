package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func Filter[TYPE any](data []TYPE, f func(TYPE) bool) []TYPE {
	var result []TYPE
	for _, d := range data {
		if f(d) {
			result = append(result, d)
		}
	}
	return result
}

func CompareVersions(version1, version2 string) (int, error) {
	v1, err := parseVersion(version1)
	if err != nil {
		return 0, err
	}
	v2, err := parseVersion(version2)
	if err != nil {
		return 0, err
	}
	return compareVersions(v1, v2), nil
}

func parseVersion(version string) ([]int, error) {
	parts := strings.Split(version, ".")
	ints := make([]int, len(parts))
	for i, part := range parts {
		var err error
		ints[i], err = strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid version part: %s", part)
		}
	}
	return ints, nil
}

func compareVersions(version1, version2 []int) int {
	for i := 0; i < len(version1) && i < len(version2); i++ {
		if version1[i] < version2[i] {
			return -1
		}
		if version1[i] > version2[i] {
			return 1
		}
	}
	if len(version1) < len(version2) {
		return -1
	}
	if len(version1) > len(version2) {
		return 1
	}
	return 0
}
