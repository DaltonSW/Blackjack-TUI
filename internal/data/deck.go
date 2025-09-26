package data

import (
	"math/rand"
)

type Deck struct {
	cards []Card
}

func NewDeck(numDecks int) *Deck {
	deck := &Deck{}
	if numDecks < 0 {
		return deck
	}

	for suit := Hearts; suit <= Spades; suit++ {
		for rank := Ace; rank <= King; rank++ {
			deck.cards = append(deck.cards, Card{Suit: suit, Rank: rank})
		}
	}

	outDeck := &Deck{}

	for range numDecks {
		outDeck.cards = append(outDeck.cards, deck.cards...)
	}
	return outDeck
}

func (d *Deck) Shuffle() {
	rand.Shuffle(len(d.cards), func(i, j int) {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	})
}

func (d *Deck) Deal() Card {
	if len(d.cards) == 0 {
		panic("cannot deal from empty deck")
	}
	card := d.cards[0]
	d.cards = d.cards[1:]
	return card
}

func (d *Deck) CardsLeft() int {
	return len(d.cards)
}
