package data

import "testing"

func TestGameRoundLifecycle(t *testing.T) {
	configs := []PlayerConfig{{Name: "Alice", Bankroll: 100}}
	game, err := NewGame(1, configs)
	if err != nil {
		t.Fatalf("unexpected error creating game: %v", err)
	}
	player := game.Players()[0]

	game.deck.cards = []Card{
		{Suit: Spades, Rank: Eight}, // player card 1
		{Suit: Clubs, Rank: Ten},    // dealer card 1
		{Suit: Hearts, Rank: Three}, // player card 2
		{Suit: Diamonds, Rank: Six}, // dealer card 2
		{Suit: Spades, Rank: Two},   // player hit
		{Suit: Hearts, Rank: Nine},  // dealer draw (causes bust)
	}

	if err := game.StartRound(map[string]int{"Alice": 10}); err != nil {
		t.Fatalf("unexpected start round error: %v", err)
	}
	if game.State() != StateDealing {
		t.Fatalf("expected game state StateDealing, got %v", game.State())
	}

	if err := game.DealInitialCards(); err != nil {
		t.Fatalf("unexpected deal initial cards error: %v", err)
	}
	if game.State() != StatePlayerAction {
		t.Fatalf("expected game state StatePlayerAction, got %v", game.State())
	}

	phand := player.ActiveHand()
	if len(phand.Cards()) != 2 {
		t.Fatalf("expected player to have 2 cards, got %d", len(phand.Cards()))
	}

	if _, err := game.Hit(player); err != nil {
		t.Fatalf("unexpected hit error: %v", err)
	}
	if len(phand.Cards()) != 3 {
		t.Fatalf("expected player hand size 3 after hit, got %d", len(phand.Cards()))
	}

	if err := game.Stand(player); err != nil {
		t.Fatalf("unexpected stand error: %v", err)
	}
	if !game.ReadyForDealer() {
		t.Fatal("expected game to be ready for dealer after player stands")
	}
	if game.State() != StateDealerAction {
		t.Fatalf("expected state to transition to StateDealerAction, got %v", game.State())
	}

	if err := game.DealerPlay(); err != nil {
		t.Fatalf("unexpected dealer play error: %v", err)
	}

	results, err := game.SettleRound()
	if err != nil {
		t.Fatalf("unexpected settle round error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected single result, got %d", len(results))
	}
	if results[0].Outcome != OutcomeWin {
		t.Fatalf("expected player to win, got outcome %v", results[0].Outcome)
	}
	if player.Bankroll() != 110 {
		t.Fatalf("expected bankroll 110 after win, got %d", player.Bankroll())
	}
	if game.State() != StateSettled {
		t.Fatalf("expected state StateSettled after settlement, got %v", game.State())
	}
}
