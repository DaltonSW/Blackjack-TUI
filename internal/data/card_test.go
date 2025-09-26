package data

import "testing"

func TestCardValue(t *testing.T) {
	tests := []struct {
		card     Card
		expected int
	}{
		{Card{Hearts, Ace}, 11},
		{Card{Hearts, Two}, 2},
		{Card{Hearts, Ten}, 10},
		{Card{Hearts, Jack}, 10},
		{Card{Hearts, Queen}, 10},
		{Card{Hearts, King}, 10},
	}

	for _, test := range tests {
		if got := test.card.Value(); got != test.expected {
			t.Errorf("Card.Value() = %d, want %d", got, test.expected)
		}
	}
}

func TestNewDeck(t *testing.T) {
	tests := []struct {
		numDecks int
		expected int
	}{
		{-1, 0},
		{0, 0},
		{1, 52},
		{2, 104},
		{4, 208},
	}

	for _, test := range tests {
		deck := NewDeck(test.numDecks)
		if got := deck.CardsLeft(); got != test.expected {
			t.Errorf("Card.Value() = %d, want %d", got, test.expected)
		}
	}
}

func TestDeckDeal(t *testing.T) {
	deck := NewDeck(1)
	initialCount := deck.CardsLeft()
	deck.Deal()

	if deck.CardsLeft() != initialCount-1 {
		t.Errorf("Expected %d cards after dealing, got %d",
			initialCount-1, deck.CardsLeft())
	}
}
