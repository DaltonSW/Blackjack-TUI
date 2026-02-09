package data

import "testing"

func TestHandSplit(t *testing.T) {
	hand := NewHand()
	hand.AddCard(Card{Suit: Spades, Rank: Eight})
	hand.AddCard(Card{Suit: Hearts, Rank: Eight})
	hand.SetBet(20)

	if !hand.CanSplit() {
		t.Fatalf("expected hand to be splittable")
	}

	splitHand, err := hand.Split()
	if err != nil {
		t.Fatalf("unexpected error splitting hand: %v", err)
	}

	if len(hand.Cards()) != 1 {
		t.Fatalf("expected original hand to retain one card, got %d", len(hand.Cards()))
	}
	if len(splitHand.Cards()) != 1 {
		t.Fatalf("expected split hand to receive one card, got %d", len(splitHand.Cards()))
	}
	if hand.Bet() != splitHand.Bet() {
		t.Fatalf("expected split hand to inherit bet %d, got %d", hand.Bet(), splitHand.Bet())
	}
}

func TestHandSplitInvalid(t *testing.T) {
	hand := NewHand()
	hand.AddCard(Card{Suit: Spades, Rank: Eight})
	hand.AddCard(Card{Suit: Hearts, Rank: Nine})

	if hand.CanSplit() {
		t.Fatal("expected hand to be unsplittable")
	}
	if _, err := hand.Split(); err == nil {
		t.Fatal("expected an error when splitting unsplittable hand")
	}
}

func TestHandDoubleDown(t *testing.T) {
	hand := NewHand()
	hand.AddCard(Card{Suit: Clubs, Rank: Five})
	hand.AddCard(Card{Suit: Diamonds, Rank: Six})

	if err := hand.DoubleDown(); err != nil {
		t.Fatalf("unexpected error doubling down: %v", err)
	}
	if !hand.IsDoubleDown() {
		t.Fatal("expected hand to be marked as double down")
	}
	if !hand.IsStanding() {
		t.Fatal("expected hand to be set to standing after double down")
	}
	if err := hand.DoubleDown(); err == nil {
		t.Fatal("expected second double down attempt to fail")
	}
}

func TestHandSoftness(t *testing.T) {
	hand := NewHand()
	hand.AddCard(Card{Suit: Hearts, Rank: Ace})
	hand.AddCard(Card{Suit: Clubs, Rank: Six})

	if !hand.IsSoft() {
		t.Fatal("expected soft hand with Ace and Six")
	}

	hand.AddCard(Card{Suit: Diamonds, Rank: Ten})

	if hand.IsSoft() {
		t.Fatal("expected hard hand after adding Ten")
	}
}
