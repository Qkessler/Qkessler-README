package github

import (
	"context"
	"math/rand"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const PER_PAGE_NUMBER int = 50
const UNCATEGORIZED string = "Uncategorized"

type LangReposAndOrder struct {
	ReposByLang map[string][]*github.Repository
	LangOrder   []string
}

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
		ListOptions: github.ListOptions{
			PerPage: PER_PAGE_NUMBER,
		},
		Type: "public",
		Sort: "updated",
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

func GetRandomRepo(repositories []*github.Repository) *github.Repository {
	if len(repositories) == 0 {
		return nil
	}
	var repo github.Repository
	visited := map[int]bool{}

    rand.Seed(time.Now().UnixNano())
	for {
		if len(visited) == len(repositories) {
			return nil
		}
		randomIndex := rand.Intn(len(repositories))
		if present, _ := visited[randomIndex]; !present {
			visited[randomIndex] = true
		}

		repo = *repositories[randomIndex]
		if repo.Fork != nil && *repo.Fork == false {
			break
		}
	}

	return &repo
}

func GetReposByLanguage(repositories []*github.Repository) (map[string][]*github.Repository, []string, error) {
	reposPerLanguage := make(map[string][]*github.Repository)
	languageOrder := []string{}

	for _, repo := range repositories {
		language := repo.Language
		if language == nil {
			continue
		}

		if repoSlice, present := reposPerLanguage[*language]; present {
			reposPerLanguage[*language] = append(repoSlice, repo)
		} else {
			reposPerLanguage[*language] = []*github.Repository{repo}
			languageOrder = append(languageOrder, *language)
		}
	}

	return reposPerLanguage, languageOrder, nil
}
