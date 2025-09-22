package main

import "github.com/charmbracelet/lipgloss/v2"

var (
	CardStyle = lipgloss.NewStyle().Width(5).Height(3).Border(lipgloss.RoundedBorder())

	DiamondColor = lipgloss.Color("#FF8800")
	ClubColor    = lipgloss.Color("#90D5FF")
	HeartColor   = lipgloss.Color("#FF3333")
	SpadeColor   = lipgloss.Color("#8877CC")
)

var (
	ClubSymbol    = "󰣎"
	SpadeSymbol   = "󰣑"
	DiamondSymbol = "󰣏"
	HeartSymbol   = "󰋑"
)

type Hand struct {
	Cards []Card
}

func (h Hand) GetTotal() int {
	return 0
}

func (h Hand) Show() string {
	return ""
}

type Card struct {
	Symbol string
	Value  int
	Hidden bool
}

func (c Card) Show() string {
	return "Card"
}
