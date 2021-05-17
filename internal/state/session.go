package state

import "github.com/AntonAverchenkov/cards-http-service/internal/game"

// Session represents a persistent connection with a client
type Session struct {
	ID   int
	Deck *game.Deck
}
