package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/google/go-github/github"
	"github.com/sashabaranov/go-openai"
	"golang.org/x/oauth2"
)

const (
	basePath = "/tmp/temp" // local directory for storing repo files
)

type FilesData struct {
	Files []string
}

func main() {

	owner := os.Getenv("REPO_OWNER")
	if owner == "" {
		log.Fatal("REPO_OWNER is not set")
	}

	repo := os.Getenv("REPO_NAME")
	if repo == "" {
		log.Fatal("REPO_NAME is not set")
	}

	apiKeys := os.Getenv("API_KEYS")
	if apiKeys == "" {
		log.Fatal("API_KEYS is not set")
	}

	accessToken := os.Getenv("GITHUB_PAT")
	if accessToken == "" {
		log.Fatal("GITHUB_PAT is not set")
	}

	// github token
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	githubClient := github.NewClient(tc)

	// openai API key
	openaiClient := openai.NewClient(apiKeys)

	// Get the latest commit
	commits, _, err := githubClient.Repositories.ListCommits(ctx, owner, repo, &github.CommitsListOptions{ListOptions: github.ListOptions{PerPage: 1}})
	if err != nil {
		log.Fatal(err)
	}
	// latestCommit := commits[0]
	// fmt.Println("latest commit:\n", latestCommit)
	latestCommitSHA := commits[0].GetSHA()

	// Get the commit object, which includes the files
	commit, _, err := githubClient.Repositories.GetCommit(ctx, owner, repo, latestCommitSHA)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range commit.Files {
		// fmt.Println(*file.RawURL)
		getContents(ctx, githubClient, *file.Filename, latestCommitSHA, owner, repo)
	}

	f := &FilesData{}
	f.printFiles(basePath)

	for _, file := range f.Files {
		githubOpenAI(ctx, file, githubClient, openaiClient, owner, repo, latestCommitSHA)
	}
}

func githubOpenAI(ctx context.Context, file string, githubClient *github.Client, openaiClient *openai.Client, owner string, repo string, latestCommitSHA string) {
	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("ReadFile error: %v\n", err)
		return
	}
	resp, err := openaiClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			// gpt-3.5-turbo
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleUser,
					Content: fmt.Sprintf(
						"Please do detailed code review of a file content and if needed propose a fix: %s",
						content),
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	comment := fmt.Sprintf("ChatGPT's review about `%s` file:\n %s", file, resp.Choices[0].Message.Content)
	fmt.Println(comment)
	githubClient.Repositories.CreateComment(ctx, owner, repo, latestCommitSHA, &github.RepositoryComment{Body: &comment})
}

func getContents(ctx context.Context, client *github.Client, path string, commitSHA string, owner string, repo string) {

	opts := &github.RepositoryContentGetOptions{Ref: commitSHA}
	rc, err := client.Repositories.DownloadContents(ctx, owner, repo, path, opts)
	if err != nil {
		log.Fatal(err)
	}
	defer rc.Close()

	body, err := io.ReadAll(rc)
	if err != nil {
		log.Fatal(err)
	}

	local := filepath.Join(basePath, path)

	// Create the directory if it does not exist
	dir := filepath.Dir(local)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = os.WriteFile(local, body, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Downloaded file:", local)
}

func (f *FilesData) printFiles(path string) {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		filePath := filepath.Join(path, file.Name())
		if file.IsDir() {
			f.printFiles(filePath)
		} else {
			f.Files = append(f.Files, filePath)
		}
	}
}
