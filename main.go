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

type Handlers struct {
	logger *log.Logger
}

// NewHandlers creates a new Handlers object with the provided logger.
func NewHandlers(logger *log.Logger) *Handlers {
	return &Handlers{
		logger: logger,
	}
}

func main() {

	logger := log.New(os.Stdout, "gh-oai ", log.LstdFlags|log.Lshortfile|log.Ltime|log.LUTC)
	h := NewHandlers(logger)

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		h.logger.Fatal("API_KEY is not set")
	}

	accessToken := os.Getenv("GITHUB_PAT")
	if accessToken == "" {
		h.logger.Fatal("GITHUB_PAT is not set")
	}

	owner := os.Getenv("REPO_OWNER")
	if owner == "" {
		h.logger.Fatal("REPO_OWNER is not set")
	}

	repo := os.Getenv("REPO_NAME")
	if repo == "" {
		h.logger.Fatal("REPO_NAME is not set")
	}

	branch := os.Getenv("BRANCH_NAME")
	var opts *github.CommitsListOptions
	if branch == "" {
		opts = &github.CommitsListOptions{ListOptions: github.ListOptions{PerPage: 1}}
	} else {
		opts = &github.CommitsListOptions{SHA: branch, ListOptions: github.ListOptions{PerPage: 1}}
	}

	// github token
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	githubClient := github.NewClient(tc)

	// openai API key
	openaiClient := openai.NewClient(apiKey)

	// Get the latest commit
	commits, _, err := githubClient.Repositories.ListCommits(ctx, owner, repo, opts)
	if err != nil {
		h.logger.Fatal(err)
	}
	// fmt.Println("latest commit:\n", commits[0])
	latestCommitSHA := commits[0].GetSHA()

	// Get the commit object, which includes the files
	commit, _, err := githubClient.Repositories.GetCommit(ctx, owner, repo, latestCommitSHA)
	if err != nil {
		h.logger.Fatal(err)
	}

	for _, file := range commit.Files {
		// fmt.Println(*file.RawURL)
		h.getContents(ctx, githubClient, *file.Filename, owner, repo, latestCommitSHA, branch)
	}

	f := &FilesData{}
	h.printFiles(f, basePath)

	for _, file := range f.Files {
		h.githubOpenAI(ctx, file, githubClient, openaiClient, owner, repo, latestCommitSHA)
	}
}

// githubOpenAI uses OpenAI's GPT-3.5-turbo model to review the content of a file,
// and creates a comment on the commit in the GitHub repository.
func (h *Handlers) githubOpenAI(ctx context.Context, file string, githubClient *github.Client, openaiClient *openai.Client, owner string, repo string, latestCommitSHA string) {
	content, err := os.ReadFile(file)
	if err != nil {
		h.logger.Printf("ReadFile error: %v\n", err)
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
		h.logger.Printf("CreateChatCompletion error: %v\n", err)
		return
	}

	comment := fmt.Sprintf("ChatGPT's review about `%s` file:\n %s", file, resp.Choices[0].Message.Content)
	fmt.Println(comment)
	githubClient.Repositories.CreateComment(ctx, owner, repo, latestCommitSHA, &github.RepositoryComment{Body: &comment})
}

// getContents downloads the content of a file from a GitHub repository and saves it locally.
func (h *Handlers) getContents(ctx context.Context, client *github.Client, path string, owner string, repo string, branch string, commitSHA string) {

	var opts *github.RepositoryContentGetOptions
	if branch == "" {
		opts = &github.RepositoryContentGetOptions{Ref: commitSHA}

	} else {
		opts = &github.RepositoryContentGetOptions{Ref: branch}
	}

	rc, err := client.Repositories.DownloadContents(ctx, owner, repo, path, opts)
	if err != nil {
		h.logger.Fatal(err)
	}
	defer rc.Close()

	body, err := io.ReadAll(rc)
	if err != nil {
		h.logger.Fatal(err)
	}

	local := filepath.Join(basePath, path)

	// Create the directory if it does not exist
	dir := filepath.Dir(local)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			h.logger.Fatal(err)
		}
	}

	err = os.WriteFile(local, body, 0644)
	if err != nil {
		h.logger.Fatal(err)
	}

	fmt.Println("Downloaded file:", local)
}

// printFiles recursively scans a directory and appends the paths of all files it contains to f.Files.
func (h *Handlers) printFiles(f *FilesData, path string) {
	files, err := os.ReadDir(path)
	if err != nil {
		h.logger.Printf("Error reading directory %s: %v", path, err)
		return
	}

	for _, file := range files {
		filePath := filepath.Join(path, file.Name())
		if file.IsDir() {
			h.printFiles(f, filePath)
		} else {
			f.Files = append(f.Files, filePath)
		}
	}
}
