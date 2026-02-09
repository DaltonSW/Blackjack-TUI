package data

import "fmt"

type GameState int

const (
	StateBetting GameState = iota
	StateDealing
	StatePlayerAction
	StateDealerAction
	StateSettled
)

type PlayerConfig struct {
	Name     string
	Bankroll int
}

type RoundResult struct {
	Player  *Player
	Hand    *Hand
	Outcome HandOutcome
}

var (
	ErrInvalidState  = fmt.Errorf("action not allowed in current game state")
	ErrUnknownPlayer = fmt.Errorf("player is not part of this game")
)

type Game struct {
	deck    *Deck
	dealer  *Dealer
	players []*Player
	state   GameState
}

func NewGame(numDecks int, configs []PlayerConfig) (*Game, error) {
	if numDecks <= 0 {
		return nil, fmt.Errorf("number of decks must be positive")
	}
	if len(configs) == 0 {
		return nil, fmt.Errorf("at least one player required")
	}
	deck := NewDeck(numDecks)
	deck.Shuffle()
	players := make([]*Player, len(configs))
	for i, cfg := range configs {
		if cfg.Bankroll <= 0 {
			return nil, fmt.Errorf("player %s must start with a positive bankroll", cfg.Name)
		}
		players[i] = NewPlayer(cfg.Name, cfg.Bankroll)
	}
	return &Game{
		deck:    deck,
		dealer:  NewDealer(),
		players: players,
		state:   StateBetting,
	}, nil
}

func (g *Game) Deck() *Deck {
	return g.deck
}

func (g *Game) Dealer() *Dealer {
	return g.dealer
}

func (g *Game) Players() []*Player {
	return g.players
}

func (g *Game) State() GameState {
	return g.state
}

func (g *Game) StartRound(bets map[string]int) error {
	if g.state != StateBetting {
		return ErrInvalidState
	}
	g.dealer.ResetForRound()
	for _, player := range g.players {
		player.ResetForRound()
		bet, ok := bets[player.Name()]
		if !ok {
			return fmt.Errorf("missing bet for player %s", player.Name())
		}
		if err := player.PlaceBet(bet); err != nil {
			return fmt.Errorf("player %s bet failed: %w", player.Name(), err)
		}
	}
	g.state = StateDealing
	return nil
}

func (g *Game) DealInitialCards() error {
	if g.state != StateDealing {
		return ErrInvalidState
	}
	for i := 0; i < 2; i++ {
		for _, player := range g.players {
			card := g.deck.Deal()
			player.ActiveHand().AddCard(card)
		}
		dealerCard := g.deck.Deal()
		g.dealer.ActiveHand().AddCard(dealerCard)
	}
	for _, player := range g.players {
		player.SetStatus(PlayerStatusActing)
	}
	g.state = StatePlayerAction
	return nil
}

func (g *Game) Hit(player *Player) (Card, error) {
	if g.state != StatePlayerAction {
		return Card{}, ErrInvalidState
	}
	if !g.containsPlayer(player) {
		return Card{}, ErrUnknownPlayer
	}
	card := g.deck.Deal()
	active := player.ActiveHand()
	if active == nil {
		return Card{}, ErrNoActiveHand
	}
	active.AddCard(card)
	if active.IsBusted() {
		active.Stand()
		player.SetStatus(PlayerStatusStanding)
	}
	return card, nil
}

func (g *Game) Stand(player *Player) error {
	if g.state != StatePlayerAction {
		return ErrInvalidState
	}
	if !g.containsPlayer(player) {
		return ErrUnknownPlayer
	}
	active := player.ActiveHand()
	if active == nil {
		return ErrNoActiveHand
	}
	active.Stand()
	if !player.MoveToNextHand() {
		player.SetStatus(PlayerStatusStanding)
	}
	return nil
}

func (g *Game) ReadyForDealer() bool {
	if g.state != StatePlayerAction {
		return false
	}
	for _, player := range g.players {
		for _, hand := range player.Hands() {
			if !hand.IsStanding() && !hand.IsBusted() {
				return false
			}
		}
	}
	g.state = StateDealerAction
	return true
}

func (g *Game) DealerPlay() error {
	if g.state != StateDealerAction {
		return ErrInvalidState
	}
	g.dealer.RevealHoleCard()
	for g.dealer.ShouldHit() {
		card := g.deck.Deal()
		g.dealer.ActiveHand().AddCard(card)
	}
	return nil
}

func (g *Game) SettleRound() ([]RoundResult, error) {
	if g.state != StateDealerAction {
		return nil, ErrInvalidState
	}
	dealerHand := g.dealer.ActiveHand()
	dealerBlackjack := dealerHand.IsBlackjack()
	dealerBust := dealerHand.IsBusted()
	dealerValue := dealerHand.Value()
	results := make([]RoundResult, 0)
	for _, player := range g.players {
		for _, hand := range player.Hands() {
			outcome := determineOutcome(hand, dealerValue, dealerBust, dealerBlackjack)
			player.Payout(hand, outcome)
			results = append(results, RoundResult{Player: player, Hand: hand, Outcome: outcome})
		}
	}
	g.state = StateSettled
	return results, nil
}

func (g *Game) PrepareNextRound() {
	if g.state != StateSettled {
		return
	}
	g.state = StateBetting
}

func determineOutcome(hand *Hand, dealerValue int, dealerBust, dealerBlackjack bool) HandOutcome {
	if hand.IsBusted() {
		return OutcomeLose
	}
	if hand.IsBlackjack() && !dealerBlackjack {
		return OutcomeBlackjack
	}
	if dealerBust {
		return OutcomeWin
	}
	if dealerBlackjack && !hand.IsBlackjack() {
		return OutcomeLose
	}
	playerValue := hand.Value()
	if playerValue > dealerValue {
		return OutcomeWin
	}
	if playerValue < dealerValue {
		return OutcomeLose
	}
	return OutcomePush
}

func (g *Game) containsPlayer(target *Player) bool {
	for _, player := range g.players {
		if player == target {
			return true
		}
	}
	return false
}
