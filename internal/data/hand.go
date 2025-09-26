package data

import "strings"

type Hand struct {
	cards []Card
}

func NewHand() *Hand {
	return &Hand{cards: make([]Card, 0)}
}

func (h *Hand) AddCard(card Card) {
	h.cards = append(h.cards, card)
}

func (h *Hand) Cards() []Card {
	return h.cards
}

func (h *Hand) Value() int {
	value := 0
	aces := 0

	for _, card := range h.cards {
		if card.Rank == Ace {
			aces++
			value += 11
		} else {
			value += card.Value()
		}
	}

	// Adjust for aces
	for aces > 0 && value > 21 {
		value -= 10
		aces--
	}

	return value
}

func (h *Hand) IsBlackjack() bool {
	return len(h.cards) == 2 && h.Value() == 21
}

func (h *Hand) IsBusted() bool {
	return h.Value() > 21
}

func (h *Hand) String() string {
	var cardStrs []string
	for _, card := range h.cards {
		cardStrs = append(cardStrs, card.String())
	}
	return strings.Join(cardStrs, " ")
}

func (h *Hand) Clear() {
	h.cards = h.cards[:0]
}
