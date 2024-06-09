package main

import (
	"io"
	"testing"
)

func TestTapeWrite(t *testing.T) {
	file, cleanFile := createTempFile(t, "12345")
	defer cleanFile()

	tape := &tape{file}

	tape.Write([]byte("abc"))

	file.Seek(0, io.SeekStart)
	newFileContents, _ := io.ReadAll(file)

	got := string(newFileContents)
	want := "abc"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
