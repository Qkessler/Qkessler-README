package cmd

import (
	"bytes"
	"embed"
	"fmt"
	"testing"
)

var nonExistentEmbed embed.FS

//go:embed embed_for_tests.txt
var existentEmbed embed.FS

func TestWriteStaticDescriptionNonExistentEmbed(t *testing.T) {
	var buffer bytes.Buffer
	err := WriteStaticDescription(&buffer, nonExistentEmbed, "NON EXISTENT")
	if err == nil {
		t.Fatalf("Should error: %s", err)
	}
}

func TestWriteStaticDescriptionExistentEmbed(t *testing.T) {
	var buffer bytes.Buffer
	err := WriteStaticDescription(&buffer, existentEmbed, EMBED_PATH)
	if err == nil {
		t.Fatalf("Should error, since we don't have the static description: %s", err)
	}
}

func TestChunks(t *testing.T) {
	slice := []string{}
	chunkSize := 2

	fmt.Println(chunks(&slice, chunkSize))
}
