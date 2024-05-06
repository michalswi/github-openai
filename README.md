# github-openai

![](https://img.shields.io/github/stars/michalswi/github-openai)
![](https://img.shields.io/github/issues/michalswi/github-openai)
![](https://img.shields.io/github/forks/michalswi/github-openai)
![](https://img.shields.io/github/last-commit/michalswi/github-openai)
![](https://img.shields.io/github/release/michalswi/github-openai)

The Go application utilizes OpenAI's ChatGPT to review all files from the latest commit and subsequently adds comments.  

Examples are [here](https://github.com/michalswi/test/commit/7de43ee5699d9bbb41a83e181e829ab157a7f3a9#comments) or [there](https://github.com/michalswi/test/commit/c5d9951c47bd230b00709ba54aa2adab735c9844#comments).

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
