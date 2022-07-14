package main

import (
	"fmt"
	"github.com/nabeelys/pr-details/gh"
	"github.com/nabeelys/pr-details/runner"
	"os"
)

var client *gh.GithubClient

func init() {
	token, isSet := os.LookupEnv("PR_GITHUB_TOKEN")
	if !isSet {
		fmt.Println("environment variable PR_GITHUB_TOKEN is not set")
		token = "" // safety
	}
	client = gh.NewGithubClient(30, token)
}

func main() {
	// without program
	args := os.Args[1:]

	if len(args) != 3 {
		fmt.Println("Invalid number of arguments")
		os.Exit(1)
	}

	owner := args[0]
	repo := args[1]
	author := args[2]

	runner.SetClient(client)
	err := runner.Run(owner, repo, author)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}

