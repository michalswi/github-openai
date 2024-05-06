# github-openai

![](https://img.shields.io/github/stars/michalswi/github-openai)
![](https://img.shields.io/github/issues/michalswi/github-openai)
![](https://img.shields.io/github/forks/michalswi/github-openai)
![](https://img.shields.io/github/last-commit/michalswi/github-openai)
![](https://img.shields.io/github/release/michalswi/github-openai)

The Go application utilizes OpenAI's ChatGPT to review all files from the latest commit and subsequently adds comments.  

Examples:
- [main](https://github.com/michalswi/github-openai-test/commit/684842f1c83edce4c0f8cd12b545ab8febf97891#comments) branch
- [develop](https://github.com/michalswi/github-openai-test/commit/3938a0d2482b325df367c824d3ded1bed8c307a9#comments) branch
- [in progress] [github-actions]() integration

You need:
- OpenAI API key
- GitHub Personal Token

```
export API_KEY=<> &&\
export GITHUB_PAT=<> &&\
export REPO_OWNER=<owner> &&\
export REPO_NAME=<repo_name>

# [optional] if branch different than the default one (main)
export BRANCH_NAME=develop

go run main.go
```
