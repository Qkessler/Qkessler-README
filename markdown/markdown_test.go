package markdown

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/Qkessler/Qkessler-README/github"
)

func TestRepoToString(t *testing.T) {
	var output bytes.Buffer
	RepoStringToWriter(
		&output,
		&github.PersonalRepo{
			Name:        "test",
			Language:    "Language",
			Description: "description",
			URL:         "URL",
		},
	)
	fmt.Println(output.String())
	if fmt.Sprintf(FORMAT_STRING, "test", "Language", "description", "URL") != output.String() {
		t.Fatalf("Format string should be applied: ")
	}
}

func TestRepoToStringNilLanguage(t *testing.T) {
	var bytes bytes.Buffer
	RepoStringToWriter(
		&bytes,
		&github.PersonalRepo{
			Name: "test",
		},
	)
	fmt.Println(bytes.String())
}
