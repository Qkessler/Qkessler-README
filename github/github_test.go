package github

import (
	"reflect"
	"testing"
)

func TestGetRandomRepoOnlyOneRepo(t *testing.T) {
	repository := PersonalRepo{
		Name:   "test",
		IsFork: false,
	}
	repos := []*PersonalRepo{
		&repository,
	}

	randomRepo := GetRandomRepo(repos)
	if !reflect.DeepEqual(repository, *randomRepo) {
		t.Fatalf("With one repo, repository should be the one chosen as random.")
	}
}

func TestGetRandomRepoNoRepos(t *testing.T) {
	repos := []*PersonalRepo{}

	randomRepo := GetRandomRepo(repos)
	if randomRepo != nil {
		t.Fatalf("Edge case: should be nil with no repos.")
	}
}

func TestGetRandomRepoAllForks(t *testing.T) {
	repository1 := PersonalRepo{
		Name:   "test",
		IsFork: true,
	}
	repository2 := PersonalRepo{
		Name:   "test2",
		IsFork: true,
	}

	repos := []*PersonalRepo{
		&repository1,
		&repository2,
	}

	randomRepo := GetRandomRepo(repos)
	if randomRepo != nil {
		t.Fatalf("With all forks, GetRandomRepo returns nil")
	}
}

func TestGetRandomRepoOneGoodOneFork(t *testing.T) {
	repository := PersonalRepo{
		Name:   "toBeChosen",
		IsFork: false,
	}

	repos := []*PersonalRepo{
		&repository,
		{
			Name:   "fork",
			IsFork: true,
		},
	}

	randomRepo := GetRandomRepo(repos)
	if !reflect.DeepEqual(*randomRepo, repository) {
		t.Fatalf("Only the non-fork can be chosen: %+vs", randomRepo)
	}
}

func TestGetReposByLanguageSameLanguage(t *testing.T) {
	repository := PersonalRepo{
		Language: "Language",
		IsFork:   false,
	}
	repository2 := PersonalRepo{
		Language: "Language",
		IsFork:   false,
	}
	repos := []*PersonalRepo{
		&repository,
		&repository2,
	}

	reposByLang, _, err := GetReposByLanguage(repos)
	if err != nil {
		t.Fatalf("Function shouldn't error with right input")
	}

	if !reflect.DeepEqual(reposByLang, map[string][]*PersonalRepo{
		"Language": {&repository, &repository2},
	}) {
		t.Fatalf("should have the right languages: %+vs", reposByLang)
	}
}

func TestGetReposByLanguageNilLanguage(t *testing.T) {
	var repository PersonalRepo
	repos := []*PersonalRepo{
		&repository,
	}

	reposByLang, _, err := GetReposByLanguage(repos)
	if err != nil {
		t.Fatalf("Function shouldn't error with right input")
	}

	if !reflect.DeepEqual(reposByLang, map[string][]*PersonalRepo{}) {
		t.Fatalf("With nil language, we should have an empty map: %+v", reposByLang)
	}
}

func TestGetReposByLanguageOrder(t *testing.T) {
	repository := PersonalRepo{
		Language: "1",
		IsFork:   false,
	}
	repository2 := PersonalRepo{
		Language: "2",
		IsFork:   false,
	}
	repos := []*PersonalRepo{
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
