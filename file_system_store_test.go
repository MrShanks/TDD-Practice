package main

import (
	"strings"
	"testing"
)

func TestFileSystemStore(t *testing.T) {

	t.Run("get league from reader", func(t *testing.T) {
		database := strings.NewReader(`[
				{"Name":"Gianni","Wins":10},
				{"Name":"Pino","Wins":45}]`)

		store := FileSystemPlayerStore{database}

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

	t.Run("get player score", func(t *testing.T) {
		database := strings.NewReader(`[
				{"Name":"Gianni","Wins":10},
				{"Name":"Pino","Wins":45}]`)

		store := FileSystemPlayerStore{database}

		got := store.GetPlayerScore("Pino")
		want := 45
		assertScoreEqual(t, got, want)
	})
}

func assertScoreEqual(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}
