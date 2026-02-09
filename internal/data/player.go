package data

import (
	"errors"
)

// PlayerStatus captures the lifecycle of a player within a single round.
type PlayerStatus int

const (
	PlayerStatusWaiting PlayerStatus = iota
	PlayerStatusActing
	PlayerStatusStanding
	PlayerStatusSettled
)

// HandOutcome models the result of a player's hand versus the dealer.
type HandOutcome int

const (
	OutcomeLose HandOutcome = iota
	OutcomePush
	OutcomeWin
	OutcomeBlackjack
)

var (
	ErrInvalidBet           = errors.New("bet must be greater than zero")
	ErrInsufficientBankroll = errors.New("insufficient bankroll to complete action")
	ErrNoActiveHand         = errors.New("no active hand available")
	ErrSplitNotAllowed      = errors.New("active hand cannot be split")
)

type Player struct {
	name     string
	bankroll int
	hands    []*Hand
	active   int
	status   PlayerStatus
}

func NewPlayer(name string, bankroll int) *Player {
	return &Player{
		name:     name,
		bankroll: bankroll,
		hands:    []*Hand{NewHand()},
		active:   0,
		status:   PlayerStatusWaiting,
	}
}

func (p *Player) Name() string {
	return p.name
}

func (p *Player) Bankroll() int {
	return p.bankroll
}

func (p *Player) Status() PlayerStatus {
	return p.status
}

func (p *Player) SetStatus(status PlayerStatus) {
	p.status = status
}

func (p *Player) Hands() []*Hand {
	return p.hands
}

func (p *Player) ActiveHand() *Hand {
	if len(p.hands) == 0 || p.active >= len(p.hands) {
		return nil
	}
	return p.hands[p.active]
}

func (p *Player) ActiveHandIndex() int {
	return p.active
}

func (p *Player) MoveToNextHand() bool {
	if len(p.hands) == 0 {
		return false
	}
	for i := p.active + 1; i < len(p.hands); i++ {
		hand := p.hands[i]
		if !hand.IsStanding() && !hand.IsBusted() {
			p.active = i
			p.status = PlayerStatusActing
			return true
		}
	}
	p.status = PlayerStatusStanding
	return false
}

func (p *Player) ResetForRound() {
	p.hands = []*Hand{NewHand()}
	p.active = 0
	p.status = PlayerStatusWaiting
}

func (p *Player) PlaceBet(amount int) error {
	if amount <= 0 {
		return ErrInvalidBet
	}
	if amount > p.bankroll {
		return ErrInsufficientBankroll
	}
	hand := p.ActiveHand()
	if hand == nil {
		return ErrNoActiveHand
	}
	p.bankroll -= amount
	hand.SetBet(amount)
	return nil
}

func (p *Player) SplitActiveHand() (*Hand, error) {
	hand := p.ActiveHand()
	if hand == nil {
		return nil, ErrNoActiveHand
	}
	if !hand.CanSplit() {
		return nil, ErrSplitNotAllowed
	}
	if hand.Bet() > p.bankroll {
		return nil, ErrInsufficientBankroll
	}
	newHand, err := hand.Split()
	if err != nil {
		return nil, err
	}
	p.bankroll -= hand.Bet()
	p.hands = append(p.hands, nil)
	copy(p.hands[p.active+2:], p.hands[p.active+1:])
	p.hands[p.active+1] = newHand
	return newHand, nil
}

func (p *Player) DoubleDownActiveHand() error {
	hand := p.ActiveHand()
	if hand == nil {
		return ErrNoActiveHand
	}
	bet := hand.Bet()
	if bet == 0 {
		return ErrInvalidBet
	}
	if bet > p.bankroll {
		return ErrInsufficientBankroll
	}
	if err := hand.DoubleDown(); err != nil {
		return err
	}
	p.bankroll -= bet
	hand.SetBet(bet * 2)
	return nil
}

func (p *Player) Payout(hand *Hand, outcome HandOutcome) {
	switch outcome {
	case OutcomeLose:
		// Bet already removed from bankroll during PlaceBet.
	case OutcomePush:
		p.bankroll += hand.Bet()
	case OutcomeWin:
		p.bankroll += hand.Bet() * 2
	case OutcomeBlackjack:
		p.bankroll += hand.Bet() + (hand.Bet()*3)/2
	}
	hand.SetBet(0)
	p.status = PlayerStatusSettled
}

type Dealer struct {
	*Player
	holeCardHidden bool
}

func NewDealer() *Dealer {
	return &Dealer{
		Player:         NewPlayer("Dealer", 0),
		holeCardHidden: true,
	}
}

func (d *Dealer) ResetForRound() {
	d.Player.ResetForRound()
	d.holeCardHidden = true
}

func (d *Dealer) ShouldHit() bool {
	hand := d.ActiveHand()
	if hand == nil {
		return false
	}
	value := hand.Value()
	if value < 17 {
		return true
	}
	return value == 17 && hand.IsSoft()
}

func (d *Dealer) ShowFirstCard() string {
	hand := d.ActiveHand()
	if hand == nil {
		return ""
	}
	cards := hand.Cards()
	if len(cards) == 0 {
		return ""
	}
	if d.holeCardHidden && len(cards) > 1 {
		return cards[0].String() + " [?]"
	}
	return hand.String()
}

func (d *Dealer) RevealHoleCard() {
	d.holeCardHidden = false
}

func (d *Dealer) HoleCardHidden() bool {
	return d.holeCardHidden
}
