package state

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/AntonAverchenkov/cards-http-service/internal/game"
	"github.com/hashicorp/go-multierror"
)

// Persist will write sessions information to the given file
func (s *SessionManager) Persist(path string) (errs error) {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("could not create %q: %w", path, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			errs = multierror.Append(errs, fmt.Errorf("could not close %q: %w", path, err))
		}
	}()

	for _, session := range s.sessions {
		if _, err := fmt.Fprintf(f, "%s %s\n", session.Id, session.Deck.Serialize()); err != nil {
			return fmt.Errorf("could not write to %q file: %w", path, err)
		}
	}

	return nil
}

// Restore will restore sessions from the given file
func Restore(path string) (_ *SessionManager, errs error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open %q: %w", path, err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			errs = multierror.Append(errs, fmt.Errorf("could not close %q: %w", path, err))
		}
	}()

	scanner := bufio.NewScanner(f)

	sessions := make(map[string]Session, 0)

	for scanner.Scan() {
		line := scanner.Text()

		tokens := strings.Split(line, " ")

		if len(tokens) != 2 {
			return nil, fmt.Errorf("%q has incorrect number of tokens", path)
		}

		deck, err := game.DeckDeserialize(tokens[1])
		if err != nil {
			return nil, fmt.Errorf("%q: deck could not be parsed: %w", path, err)
		}

		sessions[tokens[0]] = Session{
			Id:   tokens[0],
			Deck: deck,
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("could not read %q %w", path, err)
	}

	return &SessionManager{
		sessions: sessions,
	}, nil
}
