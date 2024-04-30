# [in progress]

go app (openai/chatgpt) is reviewing each file from a repo and adding comments for each file.  
Example [here](https://github.com/michalswi/test/commit/7de43ee5699d9bbb41a83e181e829ab157a7f3a9#comments) .

```
export API_KEYS=<> &&\
export GITHUB_PAT=<> &&\
export REPO_OWNER=<owner> &&\
export REPO_NAME=<repo_name>
```

#### todo
- docker image that you add to your github actions workflow to review the latest commit automatically
