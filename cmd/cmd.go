package cmd

import (
	"context"
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"

	gh "github.com/Qkessler/Qkessler-README/github"
	"github.com/Qkessler/Qkessler-README/markdown"
	"github.com/google/go-github/github"
)

const EMBED_PATH string = "assets/static-description.md"

func WriteStaticDescription(writer io.Writer, description embed.FS) error {
	text, err := description.ReadFile(EMBED_PATH)
	if err != nil {
		fmt.Println("Couldn't read embedded static description path.")
	}

	io.WriteString(writer, string(text))
	return err
}

func OpenFileAndWriteDescription(
	writeDescriptionChan *chan error,
	content embed.FS,
	filePath string,
) {
	fd, err := os.OpenFile(filePath, os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Couldn't open file for writing the description.")
	}

	*writeDescriptionChan <- WriteStaticDescription(fd, content)

	close(*writeDescriptionChan)
}

func GetRandomRepoToString(
	randomRepoChan *chan error,
	readmeFileFd io.Writer,
	svgWriter io.Writer,
	repositories []*github.Repository,
) {
	randomRepo := gh.GetRandomRepo(repositories[:10])

	err := markdown.WriteHighlightedRepoUrl(readmeFileFd, *randomRepo.HTMLURL)
	if err != nil {
		fmt.Println(err)
	}
	*randomRepoChan <- markdown.RepoStringToWriter(svgWriter, os.Stdout, randomRepo)
	close(*randomRepoChan)
}

func PullReposAndLanguageOrder(
	reposOrderChan *chan gh.LangReposAndOrder,
	repositories *[]*github.Repository,
) {
	reposByLanguage, languageOrder, err := gh.GetReposByLanguage(*repositories)
	if err != nil {
		fmt.Println(err)
		return
	}

	*reposOrderChan <- gh.LangReposAndOrder{
		ReposByLang: reposByLanguage,
		LangOrder:   languageOrder,
	}
}

func Execute(content embed.FS) {
	// To be able to use this, we should run `pass show gh-access-token`
	// to build the env variable before running this.
	accessToken := os.Getenv("GH_ACCESS_TOKEN")

	context := context.Background()

	client := gh.InitGithubClient(context, accessToken)

	readmeFile, err := filepath.Abs("Qkessler/README.md")
	if err != nil {
		fmt.Println(err)
	}

	os.Truncate(readmeFile, 0)
	writeDescriptionChan := make(chan error)
	go OpenFileAndWriteDescription(
		&writeDescriptionChan,
		content,
		readmeFile,
	)

	repositories, err := gh.PullRepositories(&context, client)
	if err != nil {
		fmt.Println("Error pulling repositories: err: ", err)
		return
	}

	descriptionError := <-writeDescriptionChan
	if descriptionError != nil {
		fmt.Println("Error writing description: ", err)
		return
	}

	svgPath, err := filepath.Abs("Qkessler/src/repo-card.svg")
	if err != nil {
		fmt.Println(err)
		return
	}

	os.Truncate(svgPath, 0)
	fd, err := os.OpenFile(svgPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Couldn't open file for random repo.")
		return
	}

	readmeFileFd, err := os.OpenFile(readmeFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Open readme file for appending failed")
	}
	randomRepoChan := make(chan error)
	go GetRandomRepoToString(&randomRepoChan, readmeFileFd, fd, repositories)

	randomRepoError := <-randomRepoChan
	if randomRepoError != nil {
		fmt.Println("Error reading the randomRepo and writing to file: ", randomRepoError)
		return
	}

	fmt.Println("Wrote description:")
	bytes, err := os.ReadFile(readmeFile)
	if err != nil {
		fmt.Println("Error reading file after writing description: ", descriptionError)
		return
	}
	fmt.Println(string(bytes))

	// reposByLanguage, order := <-reposByLanguageChan
	// fmt.Println(reposByLanguage, order)
}
