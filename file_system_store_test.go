package main

import (
	"io"
	"os"
	"testing"
)

func TestFileSystemStore(t *testing.T) {

	t.Run("get league from reader", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
				{"Name":"Gianni","Wins":10},
				{"Name":"Pino","Wins":45}]`)
		defer cleanDatabase()

		store := NewFileSystemPlayerStore(database)

		got := store.GetLeague()

		want := []Player{
			{"Gianni", 10},
			{"Pino", 45},
		}

		assertLeague(t, got, want)

		// read again
		got = store.GetLeague()
		assertLeague(t, got, want)
	})

	t.Run("record win for existing player", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
				{"Name":"Gianni","Wins":10},
				{"Name":"Pino","Wins":45}]`)
		defer cleanDatabase()

		store := NewFileSystemPlayerStore(database)

		store.RecordWin("Pino")

		got := store.GetPlayerScore("Pino")
		want := 46

		assertScoreEqual(t, got, want)
	})
	t.Run("record win for new player", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
				{"Name":"Gianni","Wins":10},
				{"Name":"Pino","Wins":45}]`)
		defer cleanDatabase()
		store := NewFileSystemPlayerStore(database)

		store.RecordWin("Marco")

		got := store.GetPlayerScore("Marco")
		want := 1

		assertScoreEqual(t, got, want)
	})

	t.Run("get player score", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
				{"Name":"Gianni","Wins":10},
				{"Name":"Pino","Wins":45}]`)
		defer cleanDatabase()

		store := NewFileSystemPlayerStore(database)

		got := store.GetPlayerScore("Pino")
		want := 45
		assertScoreEqual(t, got, want)
	})
}

func createTempFile(t testing.TB, initialData string) (io.ReadWriteSeeker, func()) {
	t.Helper()

	tmpfile, err := os.CreateTemp("", "db")
	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	tmpfile.Write([]byte(initialData))
	removeFile := func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}

	return tmpfile, removeFile
}

func assertScoreEqual(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}
