package data

import (
	"fmt"
	"strings"
)

type Hand struct {
	cards   []Card
	bet     int
	stood   bool
	doubled bool
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

func (h *Hand) Bet() int {
	return h.bet
}

func (h *Hand) SetBet(amount int) {
	h.bet = amount
}

func (h *Hand) Stand() {
	h.stood = true
}

func (h *Hand) IsStanding() bool {
	return h.stood
}

func (h *Hand) DoubleDown() error {
	if h.doubled {
		return fmt.Errorf("hand already doubled down")
	}
	if len(h.cards) != 2 {
		return fmt.Errorf("double down requires exactly two cards")
	}
	h.doubled = true
	h.stood = true
	return nil
}

func (h *Hand) IsDoubleDown() bool {
	return h.doubled
}

func (h *Hand) CanSplit() bool {
	if len(h.cards) != 2 {
		return false
	}
	first := h.cards[0]
	second := h.cards[1]
	if first.Rank == second.Rank {
		return true
	}
	// Treat all ten-value cards as splittable with each other.
	return first.Value() == 10 && second.Value() == 10
}

func (h *Hand) Split() (*Hand, error) {
	if !h.CanSplit() {
		return nil, fmt.Errorf("hand cannot be split")
	}
	second := h.cards[1]
	h.cards = h.cards[:1]
	// Reset flags for the original hand after splitting.
	h.stood = false
	h.doubled = false
	newHand := NewHand()
	newHand.cards = append(newHand.cards, second)
	newHand.bet = h.bet
	return newHand, nil
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

func (h *Hand) IsSoft() bool {
	if len(h.cards) == 0 {
		return false
	}
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
	softAces := aces
	for softAces > 0 && value > 21 {
		value -= 10
		softAces--
	}
	return softAces > 0
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
	h.bet = 0
	h.stood = false
	h.doubled = false
}
