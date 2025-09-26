package data

type Player struct {
	name string
	hand *Hand
}

func NewPlayer(name string) *Player {
	return &Player{
		name: name,
		hand: NewHand(),
	}
}

func (p *Player) Name() string {
	return p.name
}

func (p *Player) Hand() *Hand {
	return p.hand
}

func (p *Player) Reset() {
	p.hand.Clear()
}

type Dealer struct {
	*Player
}

func NewDealer() *Dealer {
	return &Dealer{
		Player: NewPlayer("Dealer"),
	}
}

func (d *Dealer) ShouldHit() bool {
	return d.hand.Value() < 17
}

func (d *Dealer) ShowFirstCard() string {
	if len(d.hand.cards) == 0 {
		return ""
	}
	return d.hand.cards[0].String() + " [?]"
}
