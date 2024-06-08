package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// StubPlayerStore is a mock player store that implements the PlayerStore interface.
// By doing this, we can write an ad hoc implementation of the playerstore that for
// example will not need a database running in order to fetch the data needed to run
// the tests. league is a slice of players that keeps track of how many wins a
// specific player has
type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
	league   []Player
}

// GetPlayerScore retrieves the player's passed score from the scores map of the
// StubPlayerStore
func (s *StubPlayerStore) GetPlayerScore(name string) int {
	return s.scores[name]
}

// RecordWin adds the winnners name to the array of the winCalls, by doing that we
// can measure the length of the winCalls array to know how many time it has been called
func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

// GetLeague returns the players stored in the league in the StubPlayerStore
func (s *StubPlayerStore) GetLeague() []Player {
	return s.league
}

func TestGETPlayers(t *testing.T) {
	store := &StubPlayerStore{
		map[string]int{
			"Gino":   20,
			"Gianni": 5,
		},
		nil,
		nil,
	}
	server := NewPlayerServer(store)

	t.Run("returns Gino's score", func(t *testing.T) {
		request := newGetScoreRequest("Gino")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, 200)
		assertResponseBody(t, response.Body.String(), "20")
	})

	t.Run("returns Gianni's score", func(t *testing.T) {
		request := newGetScoreRequest("Gianni")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, 200)
		assertResponseBody(t, response.Body.String(), "5")
	})

	t.Run("returns 404 on missing players", func(t *testing.T) {
		request := newGetScoreRequest("Pino")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Code
		want := http.StatusNotFound

		if got != want {
			t.Errorf("got status %d, want %d", got, want)
		}
	})
}

func TestStoreWins(t *testing.T) {
	store := &StubPlayerStore{
		map[string]int{},
		nil,
		nil,
	}
	server := NewPlayerServer(store)
	t.Run("records wins when POST", func(t *testing.T) {
		player := "Gianni"

		req := newPostWinRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, req)

		assertStatus(t, response.Code, http.StatusAccepted)

		if len(store.winCalls) != 1 {
			t.Errorf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
		}

		if store.winCalls[0] != player {
			t.Errorf("did not store correct winner got %q want %q", store.winCalls[0], player)
		}
	})
}

func TestLeague(t *testing.T) {
	t.Run("it returns the league table as JSON", func(t *testing.T) {
		want := []Player{
			{"Gino", 13},
			{"Marco", 20},
			{"Andrea", 13},
		}

		store := StubPlayerStore{nil, nil, want}
		server := NewPlayerServer(&store)

		request := newGetLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getLeagueFromResponse(t, response.Body)
		assertStatus(t, response.Code, http.StatusOK)
		assertLeague(t, got, want)
		assertContentType(t, response, jsonContentType)

	})
}

func getLeagueFromResponse(t testing.TB, body io.Reader) []Player {
	t.Helper()

	league, _ := NewLeague(body)
	return league
}

func assertContentType(t testing.TB, response *httptest.ResponseRecorder, want string) {
	t.Helper()

	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of %s, got %v", want, response.Result().Header)
	}
}

func assertLeague(t testing.TB, got, want []Player) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func newGetLeagueRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return req
}

func newPostWinRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func newGetScoreRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong: got %q, want %q", got, want)
	}
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("status code is wrong: got %d, want %d", got, want)
	}
}
