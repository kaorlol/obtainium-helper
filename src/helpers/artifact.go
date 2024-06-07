package helpers

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"obtainium-helper/src/utils"

	"github.com/google/go-github/v61/github"
	"golang.org/x/oauth2"
)

var (
	client *github.Client
	owner  string
	repo   string
	name   string
	branch string
)

func SetClient(token string) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token, TokenType: "Bearer"},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	client = github.NewClient(tc)
}

func getWorkflowInfo(url string) (string, string, string, string) {
	splitUrl := strings.Split(strings.Replace(url, "https://github.com/", "", 1), "/")
	owner := splitUrl[0]
	repo := splitUrl[1]
	name := strings.Split(splitUrl[4], "?")[0]
	branch := strings.Split(splitUrl[4], "branch%3A")[1]

	return owner, repo, name, branch
}

func getWorkflowLatestRun(url string) (int64, error) {
	if client == nil {
		return 0, fmt.Errorf("client not set")
	}

	owner, repo, name, branch = getWorkflowInfo(url)
	workflowRuns, _, err := client.Actions.ListWorkflowRunsByFileName(context.Background(), owner, repo, name, &github.ListWorkflowRunsOptions{
		Branch: branch,
		Status: "success",
	})

	if _, ok := err.(*github.RateLimitError); ok {
		return 0, fmt.Errorf("hit rate limit")
	}

	if len(workflowRuns.WorkflowRuns) == 0 {
		return 0, fmt.Errorf("no workflow runs found")
	}

	return workflowRuns.WorkflowRuns[0].GetID(), nil
}

func GetArtifact(url string, app utils.Download) (string, error) {
	runID, err := getWorkflowLatestRun(url)
	if err != nil {
		return "", err
	}

	artifacts, _, err := client.Actions.ListWorkflowRunArtifacts(context.Background(), owner, repo, runID, &github.ListOptions{})
	if _, ok := err.(*github.RateLimitError); ok {
		return "", fmt.Errorf("hit rate limit")
	}

	artifact := getArtifactFromPattern(app.Patterns[0], artifacts.Artifacts)
	if artifact.GetExpired() {
		return "", nil
	}

	artifactDownloadUrl, _, err := client.Actions.DownloadArtifact(context.Background(), owner, repo, artifact.GetID(), 0)
	if err != nil {
		return "", err
	}

	return artifactDownloadUrl.String(), nil
}

func getArtifactFromPattern(pattern string, artifacts []*github.Artifact) *github.Artifact {
	for _, artifact := range artifacts {
		if match, _ := regexp.MatchString(pattern, artifact.GetName()); match {
			return artifact
		}
	}

	return nil
}
