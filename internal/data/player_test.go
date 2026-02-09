package data

import "testing"

func TestPlayerPlaceBet(t *testing.T) {
	player := NewPlayer("Alice", 100)
	player.ActiveHand().AddCard(Card{Suit: Spades, Rank: Nine})
	player.ActiveHand().AddCard(Card{Suit: Clubs, Rank: Seven})

	if err := player.PlaceBet(25); err != nil {
		t.Fatalf("unexpected place bet error: %v", err)
	}
	if player.Bankroll() != 75 {
		t.Fatalf("expected bankroll 75, got %d", player.Bankroll())
	}
	if player.ActiveHand().Bet() != 25 {
		t.Fatalf("expected bet of 25 on active hand, got %d", player.ActiveHand().Bet())
	}
	if err := player.PlaceBet(-10); err == nil {
		t.Fatal("expected error for negative bet")
	}
	player.ResetForRound()
	if err := player.PlaceBet(200); err == nil {
		t.Fatal("expected error when betting more than bankroll")
	}
}

func TestPlayerSplitActiveHand(t *testing.T) {
	player := NewPlayer("Bob", 100)
	if err := player.PlaceBet(25); err != nil {
		t.Fatalf("unexpected place bet error: %v", err)
	}
	hand := player.ActiveHand()
	hand.AddCard(Card{Suit: Spades, Rank: Eight})
	hand.AddCard(Card{Suit: Hearts, Rank: Eight})

	newHand, err := player.SplitActiveHand()
	if err != nil {
		t.Fatalf("unexpected split error: %v", err)
	}
	if len(player.Hands()) != 2 {
		t.Fatalf("expected 2 hands after split, got %d", len(player.Hands()))
	}
	if player.Bankroll() != 50 {
		t.Fatalf("expected bankroll 50 after split, got %d", player.Bankroll())
	}
	if len(hand.Cards()) != 1 || len(newHand.Cards()) != 1 {
		t.Fatal("split hands should each contain a single card")
	}
	if hand.Bet() != 25 || newHand.Bet() != 25 {
		t.Fatal("split hands should retain the original bet amount")
	}
}

func TestPlayerDoubleDown(t *testing.T) {
	player := NewPlayer("Carol", 100)
	if err := player.PlaceBet(20); err != nil {
		t.Fatalf("unexpected place bet error: %v", err)
	}
	hand := player.ActiveHand()
	hand.AddCard(Card{Suit: Clubs, Rank: Nine})
	hand.AddCard(Card{Suit: Diamonds, Rank: Two})

	if err := player.DoubleDownActiveHand(); err != nil {
		t.Fatalf("unexpected double down error: %v", err)
	}
	if player.Bankroll() != 60 {
		t.Fatalf("expected bankroll 60 after double down, got %d", player.Bankroll())
	}
	if hand.Bet() != 40 {
		t.Fatalf("expected doubled bet of 40, got %d", hand.Bet())
	}
	if !hand.IsDoubleDown() {
		t.Fatal("expected hand to be flagged as double down")
	}

	player.ResetForRound()
	if err := player.PlaceBet(60); err != nil {
		t.Fatalf("unexpected place bet error after reset: %v", err)
	}
	player.ActiveHand().AddCard(Card{Suit: Clubs, Rank: Ten})
	player.ActiveHand().AddCard(Card{Suit: Hearts, Rank: Ace})
	if err := player.DoubleDownActiveHand(); err == nil {
		t.Fatal("expected double down to fail with insufficient bankroll")
	}
}

func TestPlayerPayouts(t *testing.T) {
	winPlayer := NewPlayer("Dave", 100)
	winPlayer.PlaceBet(20)
	winHand := winPlayer.ActiveHand()
	winHand.AddCard(Card{Suit: Spades, Rank: Ten})
	winHand.AddCard(Card{Suit: Hearts, Rank: Nine})
	winPlayer.Payout(winHand, OutcomeWin)
	if winPlayer.Bankroll() != 120 {
		t.Fatalf("expected bankroll 120 after win, got %d", winPlayer.Bankroll())
	}
	if winHand.Bet() != 0 {
		t.Fatal("expected bet to reset to zero after payout")
	}

	pushPlayer := NewPlayer("Eve", 100)
	pushPlayer.PlaceBet(30)
	pushHand := pushPlayer.ActiveHand()
	pushHand.AddCard(Card{Suit: Clubs, Rank: Eight})
	pushHand.AddCard(Card{Suit: Diamonds, Rank: Three})
	pushPlayer.Payout(pushHand, OutcomePush)
	if pushPlayer.Bankroll() != 100 {
		t.Fatalf("expected bankroll 100 after push, got %d", pushPlayer.Bankroll())
	}

	blackjackPlayer := NewPlayer("Frank", 100)
	blackjackPlayer.PlaceBet(40)
	blackjackHand := blackjackPlayer.ActiveHand()
	blackjackHand.AddCard(Card{Suit: Spades, Rank: Ace})
	blackjackHand.AddCard(Card{Suit: Hearts, Rank: King})
	blackjackPlayer.Payout(blackjackHand, OutcomeBlackjack)
	if blackjackPlayer.Bankroll() != 160 {
		t.Fatalf("expected bankroll 160 after blackjack, got %d", blackjackPlayer.Bankroll())
	}
}

func TestDealerBehavior(t *testing.T) {
	dealer := NewDealer()
	dealer.ActiveHand().AddCard(Card{Suit: Clubs, Rank: Ten})
	dealer.ActiveHand().AddCard(Card{Suit: Diamonds, Rank: Six})

	if !dealer.HoleCardHidden() {
		t.Fatal("expected dealer hole card to start hidden")
	}
	if dealer.ShowFirstCard() != (Card{Suit: Clubs, Rank: Ten}.String() + " [?]") {
		t.Fatalf("unexpected concealed dealer string: %s", dealer.ShowFirstCard())
	}

	dealer.RevealHoleCard()
	if dealer.HoleCardHidden() {
		t.Fatal("expected dealer hole card to be revealed")
	}
	if dealer.ShowFirstCard() == "" {
		t.Fatal("expected dealer to show full hand after reveal")
	}

	dealer.ResetForRound()
	dealer.ActiveHand().AddCard(Card{Suit: Hearts, Rank: Ace})
	dealer.ActiveHand().AddCard(Card{Suit: Spades, Rank: Six})
	if !dealer.ShouldHit() {
		t.Fatal("dealer should hit on soft 17")
	}
	dealer.ActiveHand().AddCard(Card{Suit: Clubs, Rank: Ten})
	if dealer.ShouldHit() {
		t.Fatal("dealer should stand on hard 17+")
	}
}
