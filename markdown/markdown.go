package markdown

import (
	"fmt"
	"html/template"
	"io"
	"strings"

	"github.com/google/go-github/github"
)

const FORMAT_STRING string = "%s, %s, %s, %s"
const TEMPLATE_STRING = `
<svg fill="none" viewBox="0 0 400 400" width="400" height="400" xmlns="http://www.w3.org/2000/svg">
  <foreignObject width="100%" height="100%">
    <div xmlns="http://www.w3.org/1999/xhtml">
      <style>
        .card {
			width: 300px;
			box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2);
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

func RepoStringToWriter(writer io.Writer, logWriter io.Writer, repository *github.Repository) string {
	nameOrLanguageNil := repository.Name == nil || repository.Language == nil
	descriptionOrURLNil := repository.Description == nil || repository.HTMLURL == nil
	if nameOrLanguageNil || descriptionOrURLNil {
		return ""
	}

	t, err := template.New("webpage").Parse(TEMPLATE_STRING)
	if err != nil {
		// Won't happen until we change the constant `TEMPLATE_STRING` or could be
		// a new Go version.
		fmt.Fprintf(logWriter, "Error when building template: %s", err)
	}

	data := struct {
		Name        string
		Url         string
		ImageSrc    string
		Description string
	}{
		Name:        *repository.Name,
		Url:         *repository.HTMLURL,
		ImageSrc:    "",
		Description: *repository.Description,
	}
	err = t.Execute(writer, data)

	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf(FORMAT_STRING,
		*repository.Name,
		*repository.Language,
		*repository.Description,
		*repository.HTMLURL,
	))

	return builder.String()
}
