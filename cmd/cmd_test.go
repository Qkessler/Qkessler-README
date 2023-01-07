package cmd

import (
	"bytes"
	"embed"
	"reflect"
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

func TestChunksEmptySlice(t *testing.T) {
	slice := []string{}
	chunkSize := 2

	chunks := chunks(&slice, chunkSize)

	expected := [][]string{{}}
	if !reflect.DeepEqual(chunks, expected) {
		t.Fatalf("Chunks should be with empty slice.")
	}
}

func TestJoinedSlicesEmpty(t *testing.T) {
	slice1 := []string{}
	slice2 := []string{}

	expected := []string{}
	if !reflect.DeepEqual(joinSlices(&slice1, &slice2), expected) {
		t.Fatalf("Slices should be the same, empty if no elements.")
	}
}

func TestJoinedSlicesOnly1ShouldFillWithEmpty(t *testing.T) {
	slice1 := []string{"test"}
	slice2 := []string{}

	expected := []string{"test", ""}
	joined := joinSlices(&slice1, &slice2)
	if !reflect.DeepEqual(joined, expected) {
		t.Fatalf("Slices should be the same, only one element.")
	}
}

func TestJoinedSlicesFillRightSide(t *testing.T) {
	slice1 := []string{}
	slice2 := []string{"test"}

	expected := []string{"", "test"}
	joined := joinSlices(&slice1, &slice2)
	if !reflect.DeepEqual(joined, expected) {
		t.Fatalf("Slices should be the same, only one element.")
	}
}
