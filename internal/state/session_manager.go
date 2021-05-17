package state

import "github.com/AntonAverchenkov/cards-http-service/internal/game"

// SessionManager maintains a collection of currently active sessions
type SessionManager struct {
	latest   int
	sessions map[int]Session
}

func NewSessionManaer() *SessionManager {
	return &SessionManager{
		latest:   0,
		sessions: make(map[int]Session, 0),
	}
}

func (s *SessionManager) NewSession() Session {
	s.latest++

	var (
		id      = s.latest
		session = Session{
			ID:   id,
			Deck: game.NewDeck(),
		}
	)
	s.sessions[id] = session

	return session
}

func (s *SessionManager) Get(id int) (Session, bool) {
	session, exists := s.sessions[id]

	return session, exists
}
