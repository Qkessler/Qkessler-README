package markdown

import (
	"fmt"
	"testing"

	"github.com/google/go-github/github"
)

func TestRepoToString(t *testing.T) {
	repo := RepoToString(
		&github.Repository{
			Name:     github.String("test"),
			Language: github.String("Language"),
		},
	)
	if fmt.Sprintf(FORMAT_STRING, "test", "Language") != repo {
		t.Fatalf("Format string should be applied: ")
	}
}

func TestRepoToStringNilLanguage(t *testing.T) {
	repo := RepoToString(
		&github.Repository{
			Name:     github.String("test"),
			Language: nil,
		},
	)
	if "" != repo {
		t.Fatalf("Format string should be applied: ")
	}
}
