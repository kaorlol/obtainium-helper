package utils

import (
    "archive/zip"
    "crypto/tls"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "os"
    "path/filepath"
    "regexp"
    "strings"
    "sync"
)

var (
    client = &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
            MaxIdleConns:    10,
            IdleConnTimeout: 30 * time.Second,
        },
    }
    mu sync.Mutex
)

func Request(url string, agent *string) (*http.Response, error) {
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("error creating request: %v", err)
    }

    if agent != nil {
        req.Header.Set("User-Agent", *agent)
    }

    resp, err := client.Do(req)
    if err != nil || resp.StatusCode != http.StatusOK {
        if err == nil {
            err = fmt.Errorf("unexpected status code %d", resp.StatusCode)
        }
        return nil, err
    }
    return resp, nil
}

func DownloadFile(fileURL string, app Download) (string, string, error) {
    resp, err := Request(fileURL, app.Agent)
    if err != nil {
        return "", "", err
    }
    defer resp.Body.Close()

    mu.Lock()
    defer mu.Unlock()
    if err := os.MkdirAll("out", os.ModePerm); err != nil {
        return "", "", fmt.Errorf("error creating output directory: %v", err)
    }

    identifier := extractIdentifier(fileURL, app.Identifier.Pattern)
    if identifier == app.Identifier.Latest {
        return "", "", nil
    }

    filename := extractFilename(resp)
    if strings.HasSuffix(filename, ".zip") {
        return downloadFromZip(fileURL, app)
    }

    fmt.Println("Downloading", filename, "...")
    if err := saveToFile(filepath.Join("out", filename), resp.Body); err != nil {
        return "", "", err
    }

    fmt.Println("File downloaded successfully as", filename)
    return filename, identifier, nil
}

func downloadFromZip(fileURL string, app Download) (string, string, error) {
    resp, err := Request(fileURL, app.Agent)
    if err != nil {
        return "", "", err
    }
    defer resp.Body.Close()

    zipPath := filepath.Join("out", "temp.zip")
    if err := saveToFile(zipPath, resp.Body); err != nil {
        return "", "", err
    }
    defer os.Remove(zipPath)

    r, err := zip.OpenReader(zipPath)
    if err != nil {
        return "", "", fmt.Errorf("error reading zip file: %v", err)
    }
    defer r.Close()

    biggestApk := findApk(r.File, app.Patterns[1])
    if biggestApk == nil {
        return "", "", fmt.Errorf("no apk file found in zip")
    }

    return extractAndSaveApk(biggestApk, fileURL, app.Identifier)
}

func extractAndSaveApk(apkFile *zip.File, fileURL string, Identifier Identifier) (string, string, error) {
    file, err := apkFile.Open()
    if err != nil {
        return "", "", fmt.Errorf("error opening apk file: %v", err)
    }
    defer file.Close()

    filename := apkFile.Name
    fmt.Println("Downloading", filename, "...")
    if err := saveToFile(filepath.Join("out", filename), file); err != nil {
        return "", "", err
    }

    fmt.Println("File downloaded successfully as", filename)
    return filename, extractIdentifier(fileURL, Identifier.Pattern), nil
}

func saveToFile(path string, r io.Reader) error {
    out, err := os.Create(path)
    if err != nil {
        return fmt.Errorf("error creating file: %v", err)
    }
    defer out.Close()

    if _, err := io.Copy(out, r); err != nil {
        return fmt.Errorf("error writing to file: %v", err)
    }
    return nil
}

func extractFilename(resp *http.Response) string {
    if contentDisposition := resp.Header.Get("Content-Disposition"); contentDisposition != "" {
        for _, part := range strings.Split(contentDisposition, ";") {
            if strings.HasPrefix(strings.TrimSpace(part), "filename=") {
                return strings.Trim(strings.Split(part, "=")[1], "\"")
            }
        }
    }
    return getDefaultFilename(resp.Request.URL)
}

func getDefaultFilename(u *url.URL) string {
    segments := strings.Split(u.Path, "/")
    return segments[len(segments)-1]
}

func extractIdentifier(filename, pattern string) string {
    matches := regexp.MustCompile(pattern).FindStringSubmatch(filename)
    if len(matches) > 1 {
        return matches[1]
    }
    return matches[0]
}

func findApk(files []*zip.File, pattern string) *zip.File {
    for _, f := range files {
        if strings.HasSuffix(f.Name, ".apk") && regexp.MustCompile(pattern).MatchString(f.Name) {
            return f
        }
    }
    return nil
}
