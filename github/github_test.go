package github

import (
	"reflect"
	"testing"

	"github.com/google/go-github/github"
)

func TestGetRandomRepoOnlyOneRepo(t *testing.T) {
	repository := github.Repository{
		Name: github.String("test"),
		Fork: github.Bool(false),
	}
	repos := []*github.Repository{
		&repository,
	}

	randomRepo := GetRandomRepo(repos)
	if !reflect.DeepEqual(repository, *randomRepo) {
		t.Fatalf("With one repo, repository should be the one chosen as random.")
	}
}

func TestGetRandomRepoNoRepos(t *testing.T) {
	repos := []*github.Repository{}

	randomRepo := GetRandomRepo(repos)
	if randomRepo != nil {
		t.Fatalf("Edge case: should be nil with no repos.")
	}
}

func TestGetRandomRepoAllForks(t *testing.T) {
	repository1 := github.Repository{
		Name: github.String("test"),
		Fork: github.Bool(true),
	}
	repository2 := github.Repository{
		Name: github.String("test2"),
		Fork: github.Bool(true),
	}

	repos := []*github.Repository{
		&repository1,
		&repository2,
	}

	randomRepo := GetRandomRepo(repos)
	if randomRepo != nil {
		t.Fatalf("With all forks, GetRandomRepo returns nil")
	}
}

func TestGetRandomRepoOneGoodOneFork(t *testing.T) {
	repository := github.Repository{
		Name: github.String("toBeChosen"),
		Fork: github.Bool(false),
	}

	repos := []*github.Repository{
		&repository,
		{
			Name: github.String("fork"),
			Fork: github.Bool(true),
		},
	}

	randomRepo := GetRandomRepo(repos)
	if !reflect.DeepEqual(*randomRepo, repository) {
		t.Fatalf("Only the non-fork can be chosen: %s", randomRepo)
	}
}

func TestGetReposByLanguageSameLanguage(t *testing.T) {
	repository := github.Repository{
		Language: github.String("Language"),
		Fork:     github.Bool(false),
	}
	repository2 := github.Repository{
		Language: github.String("Language"),
		Fork:     github.Bool(false),
	}
	repos := []*github.Repository{
		&repository,
		&repository2,
	}

	reposByLang, _, err := GetReposByLanguage(repos)
	if err != nil {
		t.Fatalf("Function shouldn't error with right input")
	}

	if !reflect.DeepEqual(reposByLang, map[string][]*github.Repository{
		"Language": {&repository, &repository2},
	}) {
		t.Fatalf("should have the right languages: %s", reposByLang)
	}
}

func TestGetReposByLanguageNilLanguage(t *testing.T) {
	repository := github.Repository{
		Language: nil,
	}
	repos := []*github.Repository{
		&repository,
	}

	reposByLang, _, err := GetReposByLanguage(repos)
	if err != nil {
		t.Fatalf("Function shouldn't error with right input")
	}

	if !reflect.DeepEqual(reposByLang, map[string][]*github.Repository{}) {
		t.Fatalf("With nil language, we should have an empty map: %s", reposByLang)
	}
}

func TestGetReposByLanguageOrder(t *testing.T) {
	repository := github.Repository{
		Language: github.String("1"),
		Fork: github.Bool(false),
	}
	repository2 := github.Repository{
		Language: github.String("2"),
		Fork: github.Bool(false),
	}
	repos := []*github.Repository{
		&repository,
		&repository2,
	}

	_, order, err := GetReposByLanguage(repos)
	if err != nil {
		t.Fatalf("Function shouldn't error with right input")
	}

	if !reflect.DeepEqual(order, []string{"1", "2"}) {
		t.Fatalf("Order should be kept when inserting on map: %s", order)
	}
}
