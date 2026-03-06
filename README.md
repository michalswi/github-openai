<div align="center">

# github-openai

[![stars](https://img.shields.io/github/stars/michalswi/github-openai?style=for-the-badge&color=353535)](https://github.com/michalswi/github-openai)
[![forks](https://img.shields.io/github/forks/michalswi/github-openai?style=for-the-badge&color=353535)](https://github.com/michalswi/github-openai/fork)
[![releases](https://img.shields.io/github/v/release/michalswi/github-openai?style=for-the-badge&color=353535)](https://github.com/michalswi/github-openai/releases)

the Go application utilizes OpenAI's ChatGPT (default **gpt-5-mini**) to review all files from the latest commit and subsequently adds comments.
</div>

### \# prerequisites

[OpenAI API key](https://platform.openai.com/api-keys)  
[GitHub Personal Token](https://docs.github.com/en/enterprise-server@3.9/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#about-personal-access-tokens)

### \# examples

#### **local code** review

```
export API_KEY=<> &&\
export GITHUB_PAT=<> &&\
export REPO_OWNER=<owner> &&\
export REPO_NAME=<repo_name>

# [optional] if branch different than the default one (main)
export BRANCH_NAME=develop

> I commited new changes to the specific repo:
https://github.com/michalswi/github-openai-test/commit/5f12f00872d459a1ba00438277988725ac2d6d95

$ ./github-openai_macos_arm64
Downloaded file: /tmp/temp-oJ6wj/main.go
ChatGPT's review about `main.go` file:
 Thanks — I'll review this file carefully and propose a fixed, production-ready version plus a small unit test.
 (...)
```


#### **commit** review using GitHub Actions workflow

GitHub Actions example using workflow you can find [here](https://github.com/michalswi/github-openai-test/commit/7b077c4e610bdcfa70ce44bda1facf93f5e95548) .

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
    - name: wget github-openai binary v1.0.0
      run: |
        wget https://github.com/michalswi/github-openai/releases/download/v1.0.0/github-openai_linux_amd64
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
