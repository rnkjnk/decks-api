package services

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/google/uuid"
	"github.com/rnkjnk/decks-api/internal/models"
	"github.com/rnkjnk/decks-api/internal/models/configs"
	"github.com/rnkjnk/decks-api/internal/models/dto"
)

// Decks service interface
type DecksServicer interface {
	CreateDeck(shuffle bool, cards [][2]rune) (*dto.CreateDeckResponse, error)
	OpenDeck(deckId string) (*dto.OpenDeckResponse, error)
	DrawCards(deckId string, draw uint8) (*dto.DrawCardsResponse, error)
}

type DecksService struct {
	suits     map[rune]string // Names of all possible suits
	values    map[rune]string // Names of all possible card values
	baseCards [][2]rune       // Codes of all possible cards
	decks     DecksStorer     // Repository of decks
}

func NewDecksService(config configs.DecksConfig, store DecksStorer) DecksServicer {
	// We initialize the maps of all suits and values with first character (code) as key,
	// as well as all suit and value cudes
	suits, suitCodes := dictionarize(config.Suits)
	values, valueCodes := dictionarize(config.Values)
	// now the array of all possible cards
	cards := allCards(suitCodes, valueCodes)
	// and finally, we initialize our decks service
	newDecksService := DecksService{
		suits:     suits,
		values:    values,
		baseCards: cards,
		decks:     store,
	}

	return &newDecksService
}

// Creates a new deck
func (ds *DecksService) CreateDeck(shuffle bool, cards [][2]rune) (*dto.CreateDeckResponse, error) {

	cards = intersect(ds.baseCards, cards)

	if len(cards) == 0 {
		cards = copyCards(ds.baseCards)
	}

	if shuffle {
		shuffleCards(cards)
	}

	remaining := uint8(len(cards))

	newDeck := models.Deck{
		Cards:     cards,
		Shuffled:  shuffle,
		Remaining: remaining,
	}

	_, err := ds.decks.Create(&newDeck)

	if err != nil {
		return nil, err
	}

	result := createDeckResponseFromDeck(newDeck)

	return &result, nil
}

// Opens a deck
func (ds *DecksService) OpenDeck(deckId string) (*dto.OpenDeckResponse, error) {

	id, err := uuid.Parse(deckId)
	if err != nil {
		return nil, fmt.Errorf("error parsing id: %s", deckId)
	}

	deck, err := ds.decks.Get(id)
	if err != nil {
		return nil, err
	}
	if deck.Remaining == 0 {
		return nil, fmt.Errorf("no cards remaining in deck id: %s", deckId)
	}

	remainingCards := deck.Cards[len(deck.Cards)-int(deck.Remaining):]

	result := dto.OpenDeckResponse{
		DeckId:    deck.DeckId.String(),
		Shuffled:  deck.Shuffled,
		Remaining: deck.Remaining,
		Cards:     ds.cardDtosFromIds(remainingCards),
	}

	return &result, nil
}

// Draws cards
func (ds *DecksService) DrawCards(deckId string, draw uint8) (*dto.DrawCardsResponse, error) {

	id, err := uuid.Parse(deckId)
	if err != nil {
		return nil, fmt.Errorf("error parsing id: %s", deckId)
	}

	deck, err := ds.decks.Get(id)
	if err != nil {
		return nil, err
	}
	if deck.Remaining < draw {
		return nil, fmt.Errorf("%s card(s) requested, but deck id %s has only %s card(s) left", strconv.Itoa(int(draw)), deckId, strconv.Itoa(int((deck.Remaining))))
	}

	skip := len(deck.Cards) - int(deck.Remaining)

	deck.Remaining = deck.Remaining - draw

	err = ds.decks.Put(deck)
	if err != nil {
		return nil, err
	}

	drawnCards := deck.Cards[skip : skip+int(draw)]

	result := dto.DrawCardsResponse{
		Cards: ds.cardDtosFromIds(drawnCards),
	}

	return &result, nil
}

// Function that returns all possible card codes for given arrays of suit and value codes
func allCards(suits []rune, values []rune) [][2]rune {
	// First we initialize an empty array of two-rune codes
	result := make([][2]rune, len(suits)*len(values))
	// Then we iterate over all suit and value codes to get all possible cards
	var i uint8 = 0
	for _, s := range suits {
		for _, v := range values {
			result[i][0] = v
			result[i][1] = s
			i++
		}
	}
	return result
}

// For a given array of strings, return a map with the first character of each string as the key, and an array of all keys
func dictionarize(input []string) (map[rune]string, []rune) {
	length := len(input)
	resultMap := make(map[rune]string, length)
	resultKeys := make([]rune, length)
	for i, name := range input {
		code := ([]rune(name))[0]                 // we get the first character of the string to be the code
		if _, exists := resultMap[code]; exists { // now we check if the code already exists
			panic("Bad configuration, multiple suit or value names start with same letter.")
		} else {
			resultMap[code] = name
			resultKeys[i] = code
		}
	}
	return resultMap, resultKeys
}

// Shuffles a slice of cards
func shuffleCards(cards [][2]rune) {
	// Fisher-Yates shuffle
	for i := len(cards) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		cards[i], cards[j] = cards[j], cards[i]
	}
}

// Copies a slice of cards
func copyCards(input [][2]rune) [][2]rune {
	output := make([][2]rune, len(input))
	copy(output, input)
	return output
}

// Return elements from first slice that appear in second slice
func intersect(all [][2]rune, selected [][2]rune) [][2]rune {
	result := make([][2]rune, 0, len(selected))
	for _, card := range all {
		if contains(card, selected) {
			result = append(result, card)
		}
	}
	return result
}

// Returns true if a card is found in a slice
func contains(card [2]rune, cards [][2]rune) bool {
	for _, c := range cards {
		if c == card {
			return true
		}
	}
	return false
}

func createDeckResponseFromDeck(deck models.Deck) dto.CreateDeckResponse {
	result := dto.CreateDeckResponse{
		DeckId:    deck.DeckId.String(),
		Shuffled:  deck.Shuffled,
		Remaining: deck.Remaining,
	}

	return result
}

// Returns a slice of card DTOs from a slice of IDs
func (ds *DecksService) cardDtosFromIds(cards [][2]rune) []dto.CardDto {
	result := make([]dto.CardDto, len(cards))
	for i, card := range cards {
		result[i] = dto.CardDto{
			Value: ds.values[card[0]],
			Suit:  ds.suits[card[1]],
			Code:  string(card[:]),
		}
	}
	return result
}
