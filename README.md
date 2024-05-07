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
- github-actions integration - [commit](https://github.com/michalswi/github-openai-test/commit/4467c7b457e6357ea6d7b924a66a4163b39b1301), [gh-actions](https://github.com/michalswi/github-openai-test/actions/runs/8984096714)

You need:
- [OpenAI API key](https://platform.openai.com/api-keys)
- [GitHub Personal Token](https://docs.github.com/en/enterprise-server@3.9/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#about-personal-access-tokens)

```
export API_KEY=<> &&\
export GITHUB_PAT=<> &&\
export REPO_OWNER=<owner> &&\
export REPO_NAME=<repo_name>

# [optional] if branch different than the default one (main)
export BRANCH_NAME=develop
```

#### GitHub Actions workflow example

```
name: github-openai

on:
  push:
    branches:
      - main

jobs:
  openai:
    runs-on: ubuntu-latest
    steps:
    - name: wget github-openai binary v0.2.0
      run: |
        wget https://github.com/michalswi/github-openai/releases/download/v0.2.0/github-openai_linux_amd64
        chmod +x github-openai_linux_amd64
    - name: run github-openai
      env:
        API_KEY: ${{ secrets.API_KEY }}
        GITHUB_PAT: ${{ secrets.GH_PAT }}
        REPO_OWNER: ${{ github.repository_owner }}
        REPO_NAME: ${{ github.repository }}
      run: |
        REPO_NAME=$(echo $REPO_NAME | cut -d '/' -f 2)
        # export BRANCH_NAME=develop
        ./github-openai_linux_amd64
```
