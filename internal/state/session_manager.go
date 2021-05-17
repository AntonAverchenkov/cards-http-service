package state

import (
	"crypto/rand"
	"encoding/base64"
	"io"

	"github.com/AntonAverchenkov/cards-http-service/internal/game"
)

// SessionManager maintains a collection of currently active sessions
type SessionManager struct {
	sessions map[string]Session
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]Session, 0),
	}
}

func (s *SessionManager) CreateSession() Session {
	return s.CreateSessionWith(generateUniqueSessionId())
}

func (s *SessionManager) CreateSessionWith(id string) Session {
	session := Session{
		Id:   id,
		Deck: game.NewDeck(),
	}

	s.sessions[id] = session

	return session
}

// GetOrCreateSession returns a session for the given id if it exists, creates it otherwise
func (s *SessionManager) GetOrCreateSession(id string) Session {
	session, exists := s.sessions[id]
	if exists {
		return session
	}

	return s.CreateSessionWith(id)
}

func generateUniqueSessionId() string {
	// this might be an overkill
	b := make([]byte, 32)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err)
	}

	return base64.URLEncoding.EncodeToString(b)
}
