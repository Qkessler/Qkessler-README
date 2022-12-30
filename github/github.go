package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const PER_PAGE_NUMBER int = 50

func InitGithubClient(context context.Context, accessToken string) *github.Client {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)

	tokenClient := oauth2.NewClient(context, tokenSource)
	githubClient := github.NewClient(tokenClient)

	return githubClient
}

func PullRepositories(
	context *context.Context,
	client *github.Client,
) ([]*github.Repository, error) {
	allRepos := []*github.Repository{}
	options := github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: PER_PAGE_NUMBER},
	}

	for {
		repoList, response, err := client.Repositories.List(*context, "", &options)
		if err != nil {
			return allRepos, err
		}

		allRepos = append(allRepos, repoList...)
		if response.NextPage == 0 {
			break
		}

		options.Page = response.NextPage
	}

	return allRepos, nil
}
