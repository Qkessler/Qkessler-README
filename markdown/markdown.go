package markdown

import (
	"fmt"
	"strings"

	"github.com/google/go-github/github"
)

const FORMAT_STRING string = "%s: lang: %s"

func RepoToString(repository *github.Repository) string {
	if repository.Name == nil || repository.Language == nil {
		return ""
	}

	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf(FORMAT_STRING, *repository.Name, *repository.Language))

	return builder.String()
}
