package cmd

import (
	"context"
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Qkessler/Qkessler-README/github"
	"github.com/Qkessler/Qkessler-README/markdown"
)

const EMBED_PATH string = "assets/static-description.md"
const CHUNK_SIZE int = 2

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
	repositories []*github.PersonalRepo,
) {
	randomRepo := github.GetRandomRepo(repositories[:10])

	err := markdown.WriteHighlightedRepoUrl(readmeFileFd, randomRepo.URL)
	if err != nil {
		fmt.Println(err)
	}
	*randomRepoChan <- markdown.RepoStringToWriter(svgWriter, os.Stdout, randomRepo)
	close(*randomRepoChan)
}

func PullReposAndLanguageOrder(
	reposOrderChan *chan github.LangReposAndOrder,
	repositories *[]*github.PersonalRepo,
) {
	reposByLanguage, languageOrder, err := github.GetReposByLanguage(*repositories)
	if err != nil {
		fmt.Println(err)
		return
	}

	*reposOrderChan <- github.LangReposAndOrder{
		ReposByLang: reposByLanguage,
		LangOrder:   languageOrder,
	}
	close(*reposOrderChan)
}

func Execute(content embed.FS) {
	// To be able to use this, we should run `pass show gh-access-token`
	// to build the env variable before running this.
	accessToken := os.Getenv("GH_ACCESS_TOKEN")

	context := context.Background()

	client := github.InitGithubClient(context, accessToken)

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

	repositories, err := github.PullRepositories(&context, client)
	if err != nil {
		fmt.Println("Error pulling repositories: err: ", err)
		return
	}

	svgPath, err := filepath.Abs("Qkessler/src/repo-card.svg")
	if err != nil {
		fmt.Println(err)
		return
	}

	os.Truncate(svgPath, 0)
	svgFd, err := os.OpenFile(svgPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Couldn't open file for random repo.")
		return
	}

	readmeFileFd, err := os.OpenFile(readmeFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Open readme file for appending failed")
	}
	randomRepoChan := make(chan error)
	go GetRandomRepoToString(&randomRepoChan, readmeFileFd, svgFd, repositories)

	reposByLanguageChan := make(chan github.LangReposAndOrder)
	go PullReposAndLanguageOrder(&reposByLanguageChan, &repositories)

	descriptionError := <-writeDescriptionChan
	if descriptionError != nil {
		fmt.Println("Error writing description: ", err)
		return
	}

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

	reposAndOrder := <-reposByLanguageChan
	reposByLang := reposAndOrder.ReposByLang
	for lang, repos := range reposByLang {
		fmt.Printf("Lang: %s\n", lang)
		for index, repo := range repos {
			fmt.Println(index, *repo)
		}
	}
	langOrder := reposAndOrder.LangOrder
	// for range for langOrder
	// access reposByLang for the language
	// build markdown table that we'll display on our Readme.

	// When adding the repos itselves, let's add Qkessler/Name, so we have more or
	// less the same length

	// Let's divide the number of languages in two.
	chunks := make([][]string, 0, (len(langOrder)+CHUNK_SIZE-1)/CHUNK_SIZE)
	for CHUNK_SIZE < len(langOrder) {
		langOrder, chunks = langOrder[CHUNK_SIZE:], append(chunks, langOrder[0:CHUNK_SIZE:CHUNK_SIZE])
	}
	chunks = append(chunks, langOrder)
}
