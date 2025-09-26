package data

import (
	"fmt"

	"github.com/charmbracelet/lipgloss/v2"
)

var (
	CardStyle = lipgloss.NewStyle().Width(5).Height(3).Border(lipgloss.RoundedBorder())

	DiamondColor = lipgloss.Color("#FF8800")
	ClubColor    = lipgloss.Color("#90D5FF")
	HeartColor   = lipgloss.Color("#FF3333")
	SpadeColor   = lipgloss.Color("#8877CC")
)

var (
	SuitString = []string{"♥", "♦", "♣", "♠"}
	RankString = []string{"", "A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
)

type Suit int
type Rank int

const (
	Hearts Suit = iota
	Diamonds
	Clubs
	Spades
)

const (
	Ace Rank = iota + 1
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

type Card struct {
	Suit Suit
	Rank Rank
}

func (c Card) String() string {
	return fmt.Sprintf("%s%s", RankString[c.Rank], SuitString[c.Suit])
}

func (c Card) Value() int {
	if c.Rank >= Jack {
		return 10
	}
	if c.Rank == Ace {
		return 11 // Will be adjusted in hand calculation
	}
	return int(c.Rank)
}
