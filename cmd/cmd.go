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

func WriteStaticDescription(writer io.Writer, description embed.FS, path string) error {
	text, err := description.ReadFile(path)
	if err != nil {
		fmt.Println("Couldn't read embedded static description path.")
		return err
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
		close(*writeDescriptionChan)
		return
	}

	*writeDescriptionChan <- WriteStaticDescription(fd, content, EMBED_PATH)

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
	*randomRepoChan <- markdown.RepoStringToWriter(svgWriter, randomRepo)
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
	defer svgFd.Close()
	if err != nil {
		fmt.Println("Couldn't open file for random repo.")
		return
	}

	readmeFileFd, err := os.OpenFile(readmeFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	defer readmeFileFd.Close()
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

	isHeader := true
	fmt.Fprint(readmeFileFd, "<div align='center'>\n\n")
	for _, chunk := range chunks(&langOrder, CHUNK_SIZE) {
		first, second := chunk[0], chunk[1]
		fmt.Fprintf(readmeFileFd, "|  **%s**  |  **%s**  |\n", first, second)
		if isHeader {
			fmt.Fprintf(readmeFileFd, "| :--: | :--: |\n")
			isHeader = false
		}

		firstRepos := reposByLang[first]
		secondRepos := reposByLang[second]

		joinedSlices := joinSlices(&firstRepos, &secondRepos)
		fmt.Println("JOINED: ")
		for _, v := range joinedSlices {
			fmt.Println(v)
		}
		fmt.Println("END JOINED.")
		for _, repoChunk := range chunks(&joinedSlices, CHUNK_SIZE) {
			first, second := repoChunk[0], repoChunk[1]
			markdown.RepoToStringOnTable(readmeFileFd, first, second)
		}
	}

	fmt.Fprint(readmeFileFd, "\n</div>\n")
}

func chunks[T comparable](slice *[]T, chunkSize int) [][]T {
	chunks := make([][]T, 0, (len(*slice)+chunkSize-1)/chunkSize)
	for chunkSize < len(*slice) {
		*slice, chunks = (*slice)[chunkSize:], append(chunks, (*slice)[0:chunkSize:chunkSize])
	}

	return append(chunks, *slice)
}

func joinSlices[T comparable](slice1 *[]T, slice2 *[]T) []T {
	result := []T{}
	i := 0
	j := 0

	for i < len(*slice1) && j < len(*slice2) {
		result = append(result, (*slice1)[i])
		result = append(result, (*slice2)[j])
		i += 1
		j += 1
	}

	var rest []T
	if i < len(*slice1) {
		rest = (*slice1)[i:]
	} else {
		rest = (*slice2)[j:]
	}

	for _, value := range rest {
		result = append(result, value)
		result = append(result, *new(T))
	}

	return result
}
