package markdown

import (
	"bytes"
	"errors"
	"html/template"
	"io"

	"github.com/Qkessler/Qkessler-README/github"
)

const ERROR_FIELDS_NIL string = "Name or language is Nil, or Description or URL is nil: %t, %t. Repo name: %s"
const FORMAT_STRING string = "%s, %s, %s, %s"
const TEMPLATE_STRING = `
<svg fill="none" viewBox="0 0 350 200" width="350" height="200" xmlns="http://www.w3.org/2000/svg">
  <foreignObject width="100%" height="100%">
    <div xmlns="http://www.w3.org/1999/xhtml">
		<style>
		* {
		font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol";
		}
		.github-card {
			background: white;
			display: block;
			box-sizing: border-box;
			border: 1px solid #ccc;
			margin: 20px/2;
			padding: 20px;
			color: black;
			text-decoration: none;
			font-size: 13px;
			flex: 1;
			min-width: 250px;
		}

		.github-card.h3 {
			margin-top: 0;
			color: #4078c0;
			font-size: 15px;
		}

      .github-card__meta {
	    margin-right: 20px;
	  }
	  </style>
      <div class="github-card">
		<h3>{{.Name}}</h3>
		<p>{{.Description}}</p>
        <span class="github-card__meta">
<span class="github-card__language-icon" style="color: {{.LanguageColor}};">●</span> {{.Language}}
        </span>
		<span class="github-card__meta">⭐ {{.StarNumber}}</span>
      </div>
	</div>
  </foreignObject>
</svg>`

const REPO_URL_TEMPLATE string = `
<div align="center">
    <a href="{{.Url}}">
        <img src="src/repo-card.svg" alt="Repo card which links to the Repo itself, in Github.">
    </a>
</div>

`

const REPO_ON_TABLE_TEMPLATE string = `| {{if .IsFork1 }}:small_orange_diamond:{{end}} [{{.Name1}}]({{.Url1}}) | {{if .IsFork2 }}:small_orange_diamond:{{end}} {{if .Name2}} [{{.Name2}}]({{.Url2}}) {{end}} |
`

func WriteHighlightedRepoUrl(writer io.Writer, url string) error {
	data := struct {
		Url string
	}{
		Url: url,
	}

	return WriteTemplateToWriter(writer, REPO_URL_TEMPLATE, data)
}

func WriteTemplateToWriter(writer io.Writer, templateString string, data any) error {
	templateParser, err := template.New("webpage").Parse(templateString)
	if err != nil {
		return err
	}
	var output bytes.Buffer
	err = templateParser.Execute(&output, data)
	if err != nil {
		return err
	}

	_, err = io.WriteString(writer, output.String())
	return err
}

func RepoStringToWriter(writer io.Writer, repository *github.PersonalRepo) error {
	data := struct {
		Name          string
		Url           string
		Language      string
		LanguageColor string
		Description   string
		StarNumber    int
	}{
		Name:          repository.Name,
		Url:           repository.URL,
		Language:      repository.Language,
		LanguageColor: repository.LanguageColor,
		Description:   repository.Description,
		StarNumber:    repository.StarNumber,
	}

	return WriteTemplateToWriter(writer, TEMPLATE_STRING, data)
}

func RepoToStringOnTable(
	writer io.Writer,
	repository1 *github.PersonalRepo,
	repository2 *github.PersonalRepo,
) error {
	if repository1 == nil {
		return errors.New("Can't have nil as the first repository.")
	}

	data := struct {
		Name1   string
		Url1    string
		IsFork1 bool
		Name2   string
		Url2    string
		IsFork2 bool
	}{
		Name1:   repository1.FullName,
		Url1:    repository1.URL,
		IsFork1: repository1.IsFork,
		Name2:   getOrDefault(repository2).FullName,
		Url2:    getOrDefault(repository2).URL,
		IsFork2: getOrDefault(repository2).IsFork,
	}

	return WriteTemplateToWriter(writer, REPO_ON_TABLE_TEMPLATE, data)
}

func getOrDefault[T comparable](t *T) T {
	if t == nil {
		return *new(T)
	}

	return *t
}
