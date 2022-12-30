package cmd

import (
	"context"
	"fmt"
	"os"

	gh "github.com/Qkessler/Qkessler-README/github"
	"github.com/Qkessler/Qkessler-README/markdown"
	"github.com/google/go-github/github"
)

func Execute() {
	// To be able to use this, we should run `pass show GH_ACCESS_TOKEN`
	// to build the env variable before running this.
	accessToken := os.Getenv("GH_ACCESS_TOKEN")

	context := context.Background()

	client := gh.InitGithubClient(context, accessToken)

	repositories, err := gh.PullRepositories(&context, client)
	if err != nil {
		fmt.Println("Error pulling repositories: err: ", err)
		return
	}

	randomRepoChan := make(chan string)
	go func () {
		randomRepoChan <- markdown.RepoToString(gh.GetRandomRepo(repositories[:10]))
		close(randomRepoChan)
	}()

	reposByLanguageChan := make(chan map[string][]*github.Repository)
	go func() {
		repositories, err := gh.GetReposByLanguage(repositories)
		if err != nil {
			fmt.Println(err)
			return 
		}
		reposByLanguageChan <- repositories
	}()

	randomRepo := <- randomRepoChan
	fmt.Println(randomRepo)

	reposByLanguage := <- reposByLanguageChan
	fmt.Println(reposByLanguage)
}
