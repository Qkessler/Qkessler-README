package markdown

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-github/github"
)

func TestRepoToString(t *testing.T) {
	var output bytes.Buffer
	repo := RepoStringToWriter(
		&output,
		&output,
		&github.Repository{
			Name:        github.String("test"),
			Language:    github.String("Language"),
			Description: github.String("description"),
			HTMLURL:         github.String("URL"),
		},
	)
	if fmt.Sprintf(FORMAT_STRING, "test", "Language", "description", "URL") != repo {
		t.Fatalf("Format string should be applied: ")
	}
}

func TestRepoToStringNilLanguage(t *testing.T) {
	repo := RepoStringToWriter(
		os.Stdout,
		os.Stdout,
		&github.Repository{
			Name:     github.String("test"),
			Language: nil,
		},
	)
	if "" != repo {
		t.Fatalf("Format string should be applied: ")
	}
}
