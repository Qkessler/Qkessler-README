package markdown

import (
	"bytes"
	"html/template"
	"io"

	"github.com/Qkessler/Qkessler-README/github"
)

const ERROR_FIELDS_NIL string = "Name or language is Nil, or Description or URL is nil: %t, %t. Repo name: %s"
const FORMAT_STRING string = "%s, %s, %s, %s"
const TEMPLATE_STRING = `
<svg fill="none" viewBox="0 0 400 400" width="400" height="400" xmlns="http://www.w3.org/2000/svg">
  <foreignObject width="100%" height="100%">
    <div xmlns="http://www.w3.org/1999/xhtml">
      <style>
        .card {
			width: 300px;
			box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2);
			background: white;
			border-radius: 5px;
			overflow: hidden;
			margin: 20px;
			text-align: center;
			font-family: Arial, Helvetica, sans-serif;
        }

        .card-link {
			text-decoration: none;
			color: black;
        }

        .card-header {
			display: flex;
			padding: 0px 20px 0px 10px;
			align-items: center;
        }

        .card-title {
			margin: 20px;
			font-size: 18px;
			font-weight: bold;
        }

        .card-logo {
			height: 25px;
			margin-left: 10px;
        }

        .card-description {
			margin: 20px;
			font-size: 14px;
			text-align: center;
        }
      </style>
      <div class="card">
        <a href="{{.Url}}" class="card-link">
          <div class="card-header">
            <h3 class="card-title">{{.Name}}</h3>
            <a href="" target="_blank">
              <img src="{{.ImageSrc}}" alt="Language logo" class="card-logo" />
            </a>
          </div>
          <p class="card-description">{{.Description}}</p>
        </a>
      </div>
    </div>
  </foreignObject>
</svg>`

const REPO_URL_TEMPLATE string = `
<div align="center">
    <a href="{{.Url}}">
        <img src="src/repo-card.svg" width="400" height="400" alt="Repo card which links to the Repo itself, in Github.">
    </a>
</div>
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

func RepoStringToWriter(writer io.Writer, logWriter io.Writer, repository *github.PersonalRepo) error {
	data := struct {
		Name        string
		Url         string
		ImageSrc    string
		Description string
	}{
		Name:        repository.Name,
		Url:         repository.URL,
		ImageSrc:    "",
		Description: repository.Description,
	}

	return WriteTemplateToWriter(writer, TEMPLATE_STRING, data)
}
