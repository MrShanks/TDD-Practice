package main

// This whole file is now useless and not utilized in the rest of the code, will keep it for track record but it could
// be cancelled at any time and the application would still work.
// type InMemoryPlayerStore struct {
// 	store map[string]int
// }

// func NewInMemoryPlayerStore() *InMemoryPlayerStore {
// 	return &InMemoryPlayerStore{map[string]int{}}
// }

// func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
// 	return i.store[name]
// }

// func (i *InMemoryPlayerStore) RecordWin(name string) {
// 	i.store[name]++
// }

// func (i *InMemoryPlayerStore) GetLeague() League {
// 	var league []Player
// 	for name, wins := range i.store {
// 		league = append(league, Player{name, wins})
// 	}

// 	return league
// }
