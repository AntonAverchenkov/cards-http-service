package state

import "github.com/AntonAverchenkov/cards-http-service/internal/game"

// Session represents a persistent connection with a client, each client will get their own deck
type Session struct {
	Id   string
	Deck *game.Deck
}
