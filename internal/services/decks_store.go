package services

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/rnkjnk/decks-api/internal/models"
)

// Repository interface
type DecksStorer interface {
	Get(id uuid.UUID) (*models.Deck, error)
	Put(*models.Deck) error
	Create(*models.Deck) (uuid.UUID, error)
	Delete(id uuid.UUID) error
}

// Thread safe in-memory map implementation of decks repository
type DecksInMemoryStore struct {
	// Use a mutex for safe concurrent access to the map
	mu sync.RWMutex
	// Internal map to store data
	decks map[uuid.UUID]models.Deck
}

// Creates a new in-memory repository
func NewDecksInMemoryStore() DecksStorer {
	repository := DecksInMemoryStore{
		decks: make(map[uuid.UUID]models.Deck),
	}
	return &repository
}

// Gets a deck
func (r *DecksInMemoryStore) Get(id uuid.UUID) (*models.Deck, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if data, ok := r.decks[id]; ok {
		return &data, nil
	}
	return nil, fmt.Errorf("data not found for id: %s", id)
}

// Puts a deck
func (r *DecksInMemoryStore) Put(data *models.Deck) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.decks[data.DeckId] = *data
	return nil
}

// Creates a deck
func (r *DecksInMemoryStore) Create(data *models.Deck) (uuid.UUID, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := uuid.New()
	data.DeckId = id
	r.decks[id] = *data
	return id, nil
}

// Deletes a deck
func (r *DecksInMemoryStore) Delete(id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.decks[id]; !ok {
		return fmt.Errorf("data not found for id: %s", id)
	}
	delete(r.decks, id)
	return nil
}
