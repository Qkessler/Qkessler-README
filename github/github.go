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
	ReposByLang map[string][]*PersonalRepo
	LangOrder   []string
}

type PersonalRepo struct {
	Name        string
	Description string
	Fork        bool
	Language    string
	URL         string
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
) ([]*PersonalRepo, error) {
	allRepos := []*github.Repository{}
	personalRepos := []*PersonalRepo{}
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
			return personalRepos, err
		}

		allRepos = append(allRepos, repoList...)
		if response.NextPage == 0 {
			break
		}

		options.Page = response.NextPage
	}

	for _, repo := range allRepos {
		if repo.Name == nil ||
			repo.Description == nil ||
			repo.Fork == nil ||
			repo.Language == nil ||
			repo.HTMLURL == nil {
			continue
		}

		personalRepos = append(personalRepos, &PersonalRepo{
			Name:        *repo.Name,
			Description: *repo.Description,
			Fork:        *repo.Fork,
			Language:    *repo.Language,
			URL:         *repo.HTMLURL,
		})
	}

	return personalRepos, nil
}

func GetRandomRepo(repositories []*PersonalRepo) *PersonalRepo {
	if len(repositories) == 0 {
		return nil
	}
	var repo PersonalRepo
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
		if repo.Fork == false {
			break
		}
	}

	return &repo
}

func GetReposByLanguage(repositories []*PersonalRepo) (map[string][]*PersonalRepo, []string, error) {
	reposPerLanguage := make(map[string][]*PersonalRepo)
	languageOrder := []string{}

	for _, repo := range repositories {
		language := repo.Language

		if repoSlice, present := reposPerLanguage[language]; present {
			reposPerLanguage[language] = append(repoSlice, repo)
		} else {
			reposPerLanguage[language] = []*PersonalRepo{repo}
			languageOrder = append(languageOrder, language)
		}
	}

	return reposPerLanguage, languageOrder, nil
}
