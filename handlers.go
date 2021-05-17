package main

import (
	_ "embed"
	"fmt"
	"net/http"
	"sync"

	"github.com/AntonAverchenkov/cards-http-service/internal/api"
	"github.com/AntonAverchenkov/cards-http-service/internal/game"
	"github.com/AntonAverchenkov/cards-http-service/internal/state"
	"github.com/labstack/echo/v4"
)

//go:embed doc/index.html
var documentation string

const (
	sessionCookie   = "session"
	sessionLifetime = 3600
)

type handlers struct {
	lock     sync.Mutex
	sessions *state.SessionManager
}

// (GET /) : get docuentation index.html that describes this api
func (h *handlers) Index(ctx echo.Context) error {
	return ctx.HTML(http.StatusOK, documentation)
}

// (GET /cards) : get the current state of the deck
func (h *handlers) DeckShow(ctx echo.Context) error {
	h.lock.Lock()
	defer h.lock.Unlock()

	session := h.fetchSessionSetCookie(ctx)

	return json(ctx, http.StatusOK, fromGameCards(session.Deck.Cards))
}

// (POST /cards/shuffle) : permute the deck in an unbiased way
func (h *handlers) DeckShuffle(ctx echo.Context) error {
	h.lock.Lock()
	defer h.lock.Unlock()

	session := h.fetchSessionSetCookie(ctx)
	session.Deck.Shuffle()

	return json(ctx, http.StatusOK, fromGameCards(session.Deck.Cards))
}

// (GET /cards/shuffle) : permute the deck in an unbiased way (in-browser testing helper)
func (h *handlers) DeckShuffle2(ctx echo.Context) error {
	return h.DeckShuffle(ctx)
}

// (POST /cards/deal) : deal the top card by removing it from the deck
func (h *handlers) DeckDealCard(ctx echo.Context) error {
	h.lock.Lock()
	defer h.lock.Unlock()

	session := h.fetchSessionSetCookie(ctx)

	card, err := session.Deck.DealCard()
	if err != nil {
		return json(ctx, http.StatusConflict, api.Error{Message: err.Error()})
	}

	return json(ctx, http.StatusOK, fromGameCard(card))
}

// (GET /cards/deal) : deal the top card by removing it from the deck (in-browser testing helper)
func (h *handlers) DeckDealCard2(ctx echo.Context) error {
	return h.DeckDealCard(ctx)
}

// (POST /cards/return) : return the card specified in body to the back of the deck
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

	session := h.fetchSessionSetCookie(ctx)

	err = session.Deck.ReturnCard(card)
	if err != nil {
		return json(ctx, http.StatusConflict, api.Error{Message: err.Error()})
	}

	return json(ctx, http.StatusOK, fromGameCards(session.Deck.Cards))
}

// (GET /cards/return?card={card}) : return the card specified in the '?card=' parameter to the back of the deck (in-browser testing helper)
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

	session := h.fetchSessionSetCookie(ctx)

	err = session.Deck.ReturnCard(card)
	if err != nil {
		return json(ctx, http.StatusConflict, api.Error{Message: err.Error()})
	}

	return json(ctx, http.StatusOK, fromGameCards(session.Deck.Cards))
}

// will fetch or create a new session, setting the session cookie if needed
func (h *handlers) fetchSessionSetCookie(ctx echo.Context) state.Session {

	createSessionSetCookie := func(ctx echo.Context) state.Session {
		session := h.sessions.NewSession()

		cookie := http.Cookie{
			Name:     sessionCookie,
			Value:    session.Id,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   sessionLifetime,
		}

		ctx.SetCookie(&cookie)

		return session
	}

	// check if the cookie already exists
	cookie, err := ctx.Cookie(sessionCookie)

	if err != nil || cookie.Value == "" {
		return createSessionSetCookie(ctx)
	}

	session := h.sessions.FindOrCreateSession(cookie.Value)

	return session
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
