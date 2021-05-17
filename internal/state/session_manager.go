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

func (s *SessionManager) NewSession() Session {
	var (
		id      = generateUniqueSessionId()
		session = Session{
			Id:   id,
			Deck: game.NewDeck(),
		}
	)
	s.sessions[id] = session

	return session
}

// FindOrCreateSession returns a session for the given id if it exists, returns a new session otherwise
func (s *SessionManager) FindOrCreateSession(id string) Session {
	session, exists := s.sessions[id]
	if exists {
		return session
	}

	return s.NewSession()
}

func generateUniqueSessionId() string {
	// this might be an overkill
	b := make([]byte, 32)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err)
	}

	return base64.URLEncoding.EncodeToString(b)
}
