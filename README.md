# [in progress]

go app (openai/chatgpt) is reviewing (adding comments) all files from the latest commit.  
Example [here](https://github.com/michalswi/test/commit/7de43ee5699d9bbb41a83e181e829ab157a7f3a9#comments) .

You need:
- OpenAI API key
- GitHub Personal Token

```
export API_KEY=<> &&\
export GITHUB_PAT=<> &&\
export REPO_OWNER=<owner> &&\
export REPO_NAME=<repo_name>

go run main.go
```

#### todo
- docker image that you add to your github actions workflow to review the latest commit automatically
