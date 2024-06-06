package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func GetAPKIdentifier(apkPath string) (string, error) {
	cmd := exec.Command("./aapt.exe", "dump", "badging", apkPath)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to run aapt: %v", err)
	}

	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "package:") {
			parts := strings.Split(line, " ")
			for _, part := range parts {
				if strings.HasPrefix(part, "identifierName=") {
					identifier := strings.Trim(part[len("identifierName="):], "'")
					return identifier, nil
				}
			}
		}
	}

	return "", fmt.Errorf("identifier information not found in aapt output")
}
