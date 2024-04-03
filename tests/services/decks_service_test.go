package services_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/rnkjnk/decks-api/internal/models/configs"
	"github.com/rnkjnk/decks-api/internal/models/dto"
	"github.com/rnkjnk/decks-api/internal/services"
)

// We will test with a hard-coded standard deck, so config file modifications don't spoil tests
func createMockDecksConfiguration() configs.DecksConfig {
	decksConfig := configs.DecksConfig{
		Suits:  []string{"CLUBS", "DIAMONDS", "HEARTS", "SPADES"},
		Values: []string{"ACE", "2", "3", "4", "5", "6", "7", "8", "9", "TEN", "JACK", "QUEEN", "KING"},
	}
	return decksConfig
}

func TestCreateDeck_UsesAllCardsWhenShuffled(t *testing.T) {

	service := services.NewDecksService(createMockDecksConfiguration(), services.NewDecksInMemoryStore())

	expectedResponse := &dto.CreateDeckResponse{
		DeckId:    "",
		Shuffled:  true,
		Remaining: 52,
	}

	response, err := service.CreateDeck(true, [][2]rune{})

	expectedResponse.DeckId = response.DeckId

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(response, expectedResponse) {
		t.Errorf("Unexpected response. Expected: %+v, Got: %+v", expectedResponse, response)
	}
}

func TestCreateDeck_UsesAllCardsWhenNotShuffled(t *testing.T) {

	service := services.NewDecksService(createMockDecksConfiguration(), services.NewDecksInMemoryStore())

	expectedResponse := &dto.CreateDeckResponse{
		DeckId:    "",
		Shuffled:  false,
		Remaining: 52,
	}

	response, err := service.CreateDeck(false, [][2]rune{})

	expectedResponse.DeckId = response.DeckId

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(response, expectedResponse) {
		t.Errorf("Unexpected response. Expected: %+v, Got: %+v", expectedResponse, response)
	}
}

func TestCreateDeck_ShufflesCards(t *testing.T) {

	store := services.NewDecksInMemoryStore()

	service := services.NewDecksService(createMockDecksConfiguration(), store)

	notExpectedCards := [][2]rune{
		{'A', 'C'}, {'2', 'C'}, {'3', 'C'}, {'4', 'C'}, {'5', 'C'}, {'6', 'C'}, {'7', 'C'}, {'8', 'C'}, {'9', 'C'}, {'T', 'C'}, {'J', 'C'}, {'Q', 'C'}, {'K', 'C'},
		{'A', 'D'}, {'2', 'D'}, {'3', 'D'}, {'4', 'D'}, {'5', 'D'}, {'6', 'D'}, {'7', 'D'}, {'8', 'D'}, {'9', 'D'}, {'T', 'D'}, {'J', 'D'}, {'Q', 'D'}, {'K', 'D'},
		{'A', 'H'}, {'2', 'H'}, {'3', 'H'}, {'4', 'H'}, {'5', 'H'}, {'6', 'H'}, {'7', 'H'}, {'8', 'H'}, {'9', 'H'}, {'T', 'H'}, {'J', 'H'}, {'Q', 'H'}, {'K', 'H'},
		{'A', 'S'}, {'2', 'S'}, {'3', 'S'}, {'4', 'S'}, {'5', 'S'}, {'6', 'S'}, {'7', 'S'}, {'8', 'S'}, {'9', 'S'}, {'T', 'S'}, {'J', 'S'}, {'Q', 'S'}, {'K', 'S'},
	}

	created, err := service.CreateDeck(true, [][2]rune{})

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	id, err := uuid.Parse(created.DeckId)
	if err != nil {
		t.Errorf("error parsing id: %s", created.DeckId)
	}

	deck, err := store.Get(id)
	if err != nil {
		t.Errorf("could not find deck: %s", created.DeckId)
	}

	if reflect.DeepEqual(deck.Cards, notExpectedCards) {
		t.Errorf("Unexpected response. Not expected: %+v, Got: %+v", notExpectedCards, deck.Cards)
	}
}

func TestCreateDeck_DoesNotShuffleCards(t *testing.T) {

	store := services.NewDecksInMemoryStore()

	service := services.NewDecksService(createMockDecksConfiguration(), store)

	expectedCards := [][2]rune{
		{'A', 'C'}, {'2', 'C'}, {'3', 'C'}, {'4', 'C'}, {'5', 'C'}, {'6', 'C'}, {'7', 'C'}, {'8', 'C'}, {'9', 'C'}, {'T', 'C'}, {'J', 'C'}, {'Q', 'C'}, {'K', 'C'},
		{'A', 'D'}, {'2', 'D'}, {'3', 'D'}, {'4', 'D'}, {'5', 'D'}, {'6', 'D'}, {'7', 'D'}, {'8', 'D'}, {'9', 'D'}, {'T', 'D'}, {'J', 'D'}, {'Q', 'D'}, {'K', 'D'},
		{'A', 'H'}, {'2', 'H'}, {'3', 'H'}, {'4', 'H'}, {'5', 'H'}, {'6', 'H'}, {'7', 'H'}, {'8', 'H'}, {'9', 'H'}, {'T', 'H'}, {'J', 'H'}, {'Q', 'H'}, {'K', 'H'},
		{'A', 'S'}, {'2', 'S'}, {'3', 'S'}, {'4', 'S'}, {'5', 'S'}, {'6', 'S'}, {'7', 'S'}, {'8', 'S'}, {'9', 'S'}, {'T', 'S'}, {'J', 'S'}, {'Q', 'S'}, {'K', 'S'},
	}

	created, err := service.CreateDeck(false, [][2]rune{})

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	id, err := uuid.Parse(created.DeckId)
	if err != nil {
		t.Errorf("error parsing id: %s", created.DeckId)
	}

	deck, err := store.Get(id)
	if err != nil {
		t.Errorf("could not find deck: %s", created.DeckId)
	}

	if !reflect.DeepEqual(deck.Cards, expectedCards) {
		t.Errorf("Unexpected response. Expected: %+v, Got: %+v", expectedCards, deck.Cards)
	}
}

func TestCreateDeck_IgnoresNonExistantCards(t *testing.T) {

	service := services.NewDecksService(createMockDecksConfiguration(), services.NewDecksInMemoryStore())

	expectedResponse := &dto.CreateDeckResponse{
		DeckId:    "",
		Shuffled:  false,
		Remaining: 3,
	}

	selectedCards := [][2]rune{
		{'A', 'S'},
		{'R', 'C'},
		{'2', 'L'},
		{'1', '0'},
		{'2', 'D'},
		{'3', 'H'},
	}

	response, err := service.CreateDeck(false, selectedCards)

	expectedResponse.DeckId = response.DeckId

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(response, expectedResponse) {
		t.Errorf("Unexpected response. Expected: %+v, Got: %+v", expectedResponse, response)
	}
}

func TestCreateDeck_IgnoresDuplicateCards(t *testing.T) {

	service := services.NewDecksService(createMockDecksConfiguration(), services.NewDecksInMemoryStore())

	expectedResponse := &dto.CreateDeckResponse{
		DeckId:    "",
		Shuffled:  false,
		Remaining: 1,
	}

	selectedCards := [][2]rune{
		{'A', 'S'},
		{'A', 'S'},
		{'A', 'S'},
		{'A', 'S'},
		{'A', 'S'},
	}

	response, err := service.CreateDeck(false, selectedCards)

	expectedResponse.DeckId = response.DeckId

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(response, expectedResponse) {
		t.Errorf("Unexpected response. Expected: %+v, Got: %+v", expectedResponse, response)
	}
}

func TestOpenDeck_OpensCreatedDeck(t *testing.T) {

	store := services.NewDecksInMemoryStore()

	service := services.NewDecksService(createMockDecksConfiguration(), store)

	created, err := service.CreateDeck(false, [][2]rune{
		{'A', 'C'},
		{'2', 'C'},
		{'3', 'C'},
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedResponse := dto.OpenDeckResponse{
		DeckId:    created.DeckId,
		Shuffled:  false,
		Remaining: 3,
		Cards: []dto.CardDto{
			{
				Suit:  "CLUBS",
				Value: "ACE",
				Code:  "AC",
			},
			{
				Suit:  "CLUBS",
				Value: "2",
				Code:  "2C",
			},
			{
				Suit:  "CLUBS",
				Value: "3",
				Code:  "3C",
			},
		},
	}

	response, err := service.OpenDeck(created.DeckId)
	if err != nil {
		t.Errorf("Unexpected error: %s", created.DeckId)
	}

	if !reflect.DeepEqual(*response, expectedResponse) {
		t.Errorf("Unexpected response. Expected: %+v, Got: %+v", expectedResponse, *response)
	}
}

func TestOpenDeck_ErrorIfDoesntExist(t *testing.T) {

	store := services.NewDecksInMemoryStore()

	service := services.NewDecksService(createMockDecksConfiguration(), store)

	expectedError := "data not found for id: 5b25d675-b285-4713-b976-9571a404f88a"

	_, err := service.OpenDeck("5b25d675-b285-4713-b976-9571a404f88a")
	if err == nil {
		t.Errorf("Expected error was not returned: %s", expectedError)
	}

	if expectedError != err.Error() {
		t.Errorf("Unexpected error. Expected: %+v, Got: %+v", expectedError, err.Error())
	}
}

func TestOpenDeck_ErrorIfInvalidUuid(t *testing.T) {

	store := services.NewDecksInMemoryStore()

	service := services.NewDecksService(createMockDecksConfiguration(), store)

	expectedError := "error parsing id: this is not a valid uuid"

	_, err := service.OpenDeck("this is not a valid uuid")
	if err == nil {
		t.Errorf("Expected error was not returned: %s", expectedError)
	}

	if expectedError != err.Error() {
		t.Errorf("Unexpected error. Expected: %+v, Got: %+v", expectedError, err.Error())
	}
}

func TestOpenDeck_ErrorIfNoCardsRemain(t *testing.T) {

	store := services.NewDecksInMemoryStore()

	service := services.NewDecksService(createMockDecksConfiguration(), store)

	created, err := service.CreateDeck(false, [][2]rune{
		{'A', 'C'},
		{'2', 'C'},
		{'3', 'C'},
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedError := fmt.Sprintf("no cards remaining in deck id: %s", created.DeckId)

	_, err = service.DrawCards(created.DeckId, 3)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	_, err = service.OpenDeck(created.DeckId)
	if err == nil {
		t.Errorf("Expected error was not returned: %s", expectedError)
	}

	if expectedError != err.Error() {
		t.Errorf("Unexpected error. Expected: %+v, Got: %+v", expectedError, err.Error())
	}
}

func TestDrawCards_DrawsFirstCards(t *testing.T) {

	store := services.NewDecksInMemoryStore()

	service := services.NewDecksService(createMockDecksConfiguration(), store)

	created, err := service.CreateDeck(false, [][2]rune{
		{'A', 'C'},
		{'2', 'C'},
		{'3', 'C'},
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedResponse := dto.DrawCardsResponse{
		Cards: []dto.CardDto{
			{
				Suit:  "CLUBS",
				Value: "ACE",
				Code:  "AC",
			},
			{
				Suit:  "CLUBS",
				Value: "2",
				Code:  "2C",
			},
		},
	}

	response, err := service.DrawCards(created.DeckId, 2)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(*response, expectedResponse) {
		t.Errorf("Unexpected response. Expected: %+v, Got: %+v", expectedResponse, *response)
	}
}

func TestDrawCards_DrawsSecondCards(t *testing.T) {

	store := services.NewDecksInMemoryStore()

	service := services.NewDecksService(createMockDecksConfiguration(), store)

	created, err := service.CreateDeck(false, [][2]rune{
		{'A', 'C'},
		{'2', 'C'},
		{'3', 'C'},
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedResponse := dto.DrawCardsResponse{
		Cards: []dto.CardDto{
			{
				Suit:  "CLUBS",
				Value: "2",
				Code:  "2C",
			},
			{
				Suit:  "CLUBS",
				Value: "3",
				Code:  "3C",
			},
		},
	}

	_, err = service.DrawCards(created.DeckId, 1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	response, err := service.DrawCards(created.DeckId, 2)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(*response, expectedResponse) {
		t.Errorf("Unexpected response. Expected: %+v, Got: %+v", expectedResponse, *response)
	}
}

func TestDrawCards_DecreasesRemainingCount(t *testing.T) {

	store := services.NewDecksInMemoryStore()

	service := services.NewDecksService(createMockDecksConfiguration(), store)

	created, err := service.CreateDeck(false, [][2]rune{
		{'A', 'C'},
		{'2', 'C'},
		{'3', 'C'},
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedRemainingCount := 1

	_, err = service.DrawCards(created.DeckId, 2)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	deck, err := service.OpenDeck(created.DeckId)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if deck.Remaining != uint8(expectedRemainingCount) {
		t.Errorf("Unexpected response. Expected: %+v, Got: %+v", expectedRemainingCount, deck.Remaining)
	}
}

func TestDrawCards_ErrorWhenNoCardsLeft(t *testing.T) {

	store := services.NewDecksInMemoryStore()

	service := services.NewDecksService(createMockDecksConfiguration(), store)

	created, err := service.CreateDeck(false, [][2]rune{
		{'A', 'C'},
		{'2', 'C'},
		{'3', 'C'},
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedError := fmt.Sprintf("1 card(s) requested, but deck id %s has only 0 card(s) left", created.DeckId)

	_, err = service.DrawCards(created.DeckId, 3)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	_, err = service.DrawCards(created.DeckId, 1)
	if err == nil {
		t.Errorf("Expected error was not returned: %s", expectedError)
	}

	if expectedError != err.Error() {
		t.Errorf("Unexpected error. Expected: %+v, Got: %+v", expectedError, err.Error())
	}
}

func TestDrawCards_ErrorIfDeckDoesntExist(t *testing.T) {

	store := services.NewDecksInMemoryStore()

	service := services.NewDecksService(createMockDecksConfiguration(), store)

	expectedError := "data not found for id: 5b25d675-b285-4713-b976-9571a404f88a"

	_, err := service.DrawCards("5b25d675-b285-4713-b976-9571a404f88a", 1)
	if err == nil {
		t.Errorf("Expected error was not returned: %s", expectedError)
	}

	if expectedError != err.Error() {
		t.Errorf("Unexpected error. Expected: %+v, Got: %+v", expectedError, err.Error())
	}
}

func TestDrawCards_ErrorIfInvalidDeckUuid(t *testing.T) {

	store := services.NewDecksInMemoryStore()

	service := services.NewDecksService(createMockDecksConfiguration(), store)

	expectedError := "error parsing id: this is not a valid uuid"

	_, err := service.DrawCards("this is not a valid uuid", 1)
	if err == nil {
		t.Errorf("Expected error was not returned: %s", expectedError)
	}

	if expectedError != err.Error() {
		t.Errorf("Unexpected error. Expected: %+v, Got: %+v", expectedError, err.Error())
	}
}
