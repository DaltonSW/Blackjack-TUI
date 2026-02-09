package main

import (
	"log"

	"blackjack/internal/data"
	"blackjack/internal/tui"
	tea "github.com/charmbracelet/bubbletea/v2"
)

func main() {
	game, err := data.NewGame(6, []data.PlayerConfig{{Name: "You", Bankroll: 500}})
	if err != nil {
		log.Fatalf("failed to initialize game: %v", err)
	}

	program := tea.NewProgram(tui.New(game), tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		log.Fatalf("error running TUI: %v", err)
	}
}
