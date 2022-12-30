package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/Qkessler/Qkessler-README/github"
)

func Execute() {
	// To be able to use this, we should run `pass show GH_ACCESS_TOKEN`
	// to build the env variable before running this.
	accessToken := os.Getenv("GH_ACCESS_TOKEN")

	fmt.Println(accessToken)
	context := context.Background()

	client := github.InitGithubClient(context, accessToken)

	repositories, err := github.PullRepositories(context, client)
	if err != nil {
		fmt.Println("Error pulling repositories: err: ", err)
		return
	}
	fmt.Println(repositories)
}
