package main

import (
	_ "embed"
	"fmt"
	"net/http"
	"sync"

	"github.com/AntonAverchenkov/cards-http-service/internal/api"
	"github.com/AntonAverchenkov/cards-http-service/internal/game"
	"github.com/labstack/echo/v4"
)

//go:embed doc/index.html
var documentation string

type handlers struct {
	lock sync.Mutex
	deck *game.Deck
}

// (GET /) index.html that describes this api
func (h *handlers) Index(ctx echo.Context) error {
	return ctx.HTML(http.StatusOK, documentation)
}

// (GET /cards) Get the current state of the deck
func (h *handlers) DeckShow(ctx echo.Context) error {
	h.lock.Lock()
	defer h.lock.Unlock()

	return json(ctx, http.StatusOK, fromGameCards(h.deck.Cards))
}

// (POST /cards/shuffle) Permute the deck in an unbiased way
func (h *handlers) DeckShuffle(ctx echo.Context) error {
	h.lock.Lock()
	defer h.lock.Unlock()

	h.deck.Shuffle()

	return json(ctx, http.StatusOK, fromGameCards(h.deck.Cards))
}

// (GET /cards/shuffle) Permute the deck in an unbiased way
func (h *handlers) DeckShuffle2(ctx echo.Context) error {
	return h.DeckShuffle(ctx)
}

// (POST /cards/deal) Deals the top card by removing it from the deck
func (h *handlers) DeckDealCard(ctx echo.Context) error {
	h.lock.Lock()
	defer h.lock.Unlock()

	card, err := h.deck.DealCard()
	if err != nil {
		return json(ctx, http.StatusConflict, api.Error{Message: err.Error()})
	}

	return json(ctx, http.StatusOK, fromGameCard(card))
}

// (GET /cards/deal) Deals the top card by removing it from the deck
func (h *handlers) DeckDealCard2(ctx echo.Context) error {
	return h.DeckDealCard(ctx)
}

// (POST /cards/return) Return the card specified in body to the back of the deck
func (h *handlers) DeckReturnCard(ctx echo.Context) error {
	// We expect an api.Card object in the request body
	var c api.Card
	err := ctx.Bind(&c)
	if err != nil {
		return json(ctx, http.StatusBadRequest, api.Error{Message: err.Error()})
	}

	card, err := toGameCard(c)
	if err != nil {
		return json(ctx, http.StatusBadRequest, api.Error{Message: err.Error()})
	}

	h.lock.Lock()
	defer h.lock.Unlock()

	err = h.deck.ReturnCard(card)
	if err != nil {
		return json(ctx, http.StatusConflict, api.Error{Message: err.Error()})
	}

	return json(ctx, http.StatusOK, fromGameCards(h.deck.Cards))
}

// (GET /cards/return?card={card}) Return the card specified in url parameter to the back of the deck
func (h *handlers) DeckReturnCard2(ctx echo.Context, params api.DeckReturnCard2Params) error {
	if params.Card == nil {
		return json(ctx, http.StatusBadRequest, api.Error{Message: "the required url parameter 'card' is missing"})
	}

	card, err := game.ParseCard(*params.Card)
	if err != nil {
		return json(ctx, http.StatusBadRequest, api.Error{Message: err.Error()})
	}

	h.lock.Lock()
	defer h.lock.Unlock()

	err = h.deck.ReturnCard(card)
	if err != nil {
		return json(ctx, http.StatusConflict, api.Error{Message: err.Error()})
	}

	return json(ctx, http.StatusOK, fromGameCards(h.deck.Cards))
}

// json is a formatting helper
func json(ctx echo.Context, code int, i interface{}) error {
	return ctx.JSONPretty(code, i, "  ")
}

///
/// The translation helpers below are the unfortunate side effect of dealing with swagger middle layer
///

func fromGameCard(card game.Card) api.Card {
	return api.Card{
		Value: card.Value.String(),
		Suit:  card.Suit.String(),
	}
}

func fromGameCards(cards []game.Card) []api.Card {
	var result []api.Card

	for _, card := range cards {
		result = append(result, fromGameCard(card))
	}

	return result
}

func toGameCard(card api.Card) (game.Card, error) {
	v, err := game.ParseValue(card.Value)
	if err != nil {
		return game.Card{}, fmt.Errorf("error parsing %v: %w", card, err)
	}

	s, err := game.ParseSuit(card.Suit)
	if err != nil {
		return game.Card{}, fmt.Errorf("error parsing %v: %w", card, err)
	}

	return game.Card{Value: v, Suit: s}, nil
}
