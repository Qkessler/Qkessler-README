package cmd

import (
	"context"
	"fmt"
	"os"

	gh "github.com/Qkessler/Qkessler-README/github"
	"github.com/Qkessler/Qkessler-README/markdown"
)

func Execute() {
	// To be able to use this, we should run `pass show gh-access-token`
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
	go func() {
		randomRepo := gh.GetRandomRepo(repositories[:10])
		randomRepoChan <- markdown.RepoStringToWriter(os.Stdout, os.Stdout, randomRepo)
		close(randomRepoChan)
	}()

	reposByLanguageChan := make(chan gh.LangReposAndOrder)
	go func() {
		repositories, languageOrder, err := gh.GetReposByLanguage(repositories)
		if err != nil {
			fmt.Println(err)
			return
		}
		reposByLanguageChan <- gh.LangReposAndOrder{
			ReposByLang: repositories,
			LangOrder:   languageOrder,
		}
	}()

	randomRepo := <-randomRepoChan
	fmt.Println(randomRepo)

	reposByLanguage, order := <-reposByLanguageChan
	fmt.Println(reposByLanguage, order)
}
