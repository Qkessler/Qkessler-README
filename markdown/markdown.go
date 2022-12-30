package markdown

import (
	"fmt"
	"strings"

	"github.com/google/go-github/github"
)

const FORMAT_STRING string = "%s, %s, %s, %s"

func RepoToString(repository *github.Repository) string {
	nameOrLanguageNil := repository.Name == nil || repository.Language == nil
	descriptionOrURLNil := repository.Description == nil || repository.HTMLURL == nil
	if nameOrLanguageNil || descriptionOrURLNil {
		return ""
	}

	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf(FORMAT_STRING,
		*repository.Name,
		*repository.Language,
		*repository.Description,
		*repository.HTMLURL,
	))

	return builder.String()
}
