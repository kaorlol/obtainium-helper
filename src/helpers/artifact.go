package helpers

import (
	"context"
	"fmt"
	"strings"

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

func GetArtifacts(url string, lastRunId string) (string, error) {
	runID, err := getWorkflowLatestRun(url)
	if err != nil {
		return "", err
	}

	if lastRunId == fmt.Sprintf("%d", runID) {
		return "", nil
	}

	artifacts, _, err := client.Actions.ListWorkflowRunArtifacts(context.Background(), owner, repo, runID, &github.ListOptions{
		PerPage: 1,
	})

	if _, ok := err.(*github.RateLimitError); ok {
		return "", fmt.Errorf("hit rate limit")
	}

	artifact := artifacts.Artifacts[0]
	if artifact.GetExpired() {
		return "", nil
	}

	artifactDownloadUrl, _, err := client.Actions.DownloadArtifact(context.Background(), owner, repo, artifact.GetID(), 0)
	if err != nil {
		return "", err
	}

	return artifactDownloadUrl.String(), nil
}

