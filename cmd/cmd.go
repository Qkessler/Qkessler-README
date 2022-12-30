package cmd

import (
	"context"
	"embed"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	gh "github.com/Qkessler/Qkessler-README/github"
	"github.com/Qkessler/Qkessler-README/markdown"
)

const EMBED_PATH string = "assets/static-description.md"

func WriteStaticDescription(writer io.Writer, description embed.FS) error {
	text, err := description.ReadFile(EMBED_PATH)
	if err != nil {
		fmt.Println("Couldn't read embedded static description path.")
	}

	_, err = fmt.Fprintf(writer, "%s", text)
	return err
}

func Execute(content embed.FS) {
	// To be able to use this, we should run `pass show gh-access-token`
	// to build the env variable before running this.
	accessToken := os.Getenv("GH_ACCESS_TOKEN")

	context := context.Background()

	client := gh.InitGithubClient(context, accessToken)

	repositories, err := gh.PullRepositories(&context, client)
	if err != nil {
		fmt.Println("Error pulling repositories: err: ", err)
		return
	}

	file, err := ioutil.TempFile("", "README")
	if err != nil {
		fmt.Println("Error creating Temp file for writing.")
	}
	defer os.Remove(file.Name())

	writeDescriptionChan := make(chan error)
	go func() {
		fmt.Println("file name: ", file.Name())
		fd, err := os.OpenFile(file.Name(), os.O_WRONLY, 0666)
		if err != nil {
			fmt.Println("Couldn't open file for writing the description.")
		}

		writeDescriptionChan <- WriteStaticDescription(fd, content)

		close(writeDescriptionChan)
	}()

	randomRepoChan := make(chan string)
	go func() {
		randomRepo := gh.GetRandomRepo(repositories[:10])
		randomRepoChan <- markdown.RepoStringToWriter(os.Stdout, os.Stdout, randomRepo)
		close(randomRepoChan)
	}()

	reposByLanguageChan := make(chan gh.LangReposAndOrder)
	go func() {
		repositories, languageOrder, err := gh.GetReposByLanguage(repositories)
		if err != nil {
			fmt.Println(err)
			return
		}
		reposByLanguageChan <- gh.LangReposAndOrder{
			ReposByLang: repositories,
			LangOrder:   languageOrder,
		}
	}()

	descriptionError := <-writeDescriptionChan
	if descriptionError != nil {
		fmt.Println("Error writing description: ", err)
	}
	fmt.Println("Wrote description:")
	bytes, err := os.ReadFile(file.Name())
	if err != nil {
		fmt.Println("Error reading file after writing description.")
	}
	fmt.Println(string(bytes))

	randomRepo := <-randomRepoChan
	fmt.Println(randomRepo)

	reposByLanguage, order := <-reposByLanguageChan
	fmt.Println(reposByLanguage, order)
}
