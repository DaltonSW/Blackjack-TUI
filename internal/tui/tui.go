package tui

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
	"unicode/utf8"

	"blackjack/internal/data"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

var (
	headerStyle       = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FACC15")).Background(lipgloss.Color("#1F2937")).Padding(0, 1)
	infoStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#A5B4FC"))
	sectionTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F9A8D4"))
	handBoxStyle      = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#4B5563")).Padding(0, 1).MarginTop(1)
	activeHandStyle   = handBoxStyle.Copy().BorderForeground(lipgloss.Color("#F97316"))
	valueStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#34D399")).Bold(true)
	tagStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("#F87171")).Bold(true)
	promptStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#FDE68A")).Padding(0, 1)
	inputStyle        = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#FDE68A")).Padding(0, 1)
	hotkeyBarStyle    = lipgloss.NewStyle().MarginTop(1)
	hotkeyKeyStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#38BDF8"))
	hotkeyDisabledKey = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#6B7280"))
	hotkeyLabelStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#E5E7EB"))
	messageBoxStyle   = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("#6B7280")).Padding(0, 1).MarginTop(1)
	errorStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#F87171")).Bold(true)
	cardStyle         = data.CardStyle.Copy().Width(9).Height(5)
	faceDownStyle     = cardStyle.Copy().BorderForeground(lipgloss.Color("#6B7280")).Foreground(lipgloss.Color("#6B7280"))
)

type Model struct {
	game     *data.Game
	player   *data.Player
	input    string
	messages []string
	results  []data.RoundResult
	prompt   string
	err      error
	quitting bool
}

func New(game *data.Game) *Model {
	players := game.Players()
	var player *data.Player
	if len(players) > 0 {
		player = players[0]
	}
	m := &Model{
		game:     game,
		player:   player,
		messages: []string{"Welcome to Blackjack. Place your opening bet."},
	}
	m.updatePrompt()
	return m
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch keyMsg := msg.(type) {
	case tea.KeyPressMsg:
		key := keyMsg.Key()
		text := strings.ToLower(key.Text)

		if key.Mod&tea.ModCtrl != 0 && (key.Code == 'c' || key.Code == 'C') {
			m.quitting = true
			return m, tea.Quit
		}

		if text == "q" {
			m.quitting = true
			return m, tea.Quit
		}

		switch m.game.State() {
		case data.StateBetting, data.StateSettled:
			switch key.Code {
			case tea.KeyEnter:
				trimmed := strings.TrimSpace(m.input)
				if trimmed == "" {
					break
				}
				if err := m.handleCommand(trimmed); err != nil {
					m.err = err
				} else {
					m.err = nil
				}
				m.input = ""
			case tea.KeyBackspace, tea.KeyDelete:
				m.input = trimLastRune(m.input)
			default:
				if text == "?" {
					m.showHelp()
					break
				}
				if key.Text != "" {
					r, _ := utf8.DecodeRuneInString(key.Text)
					if r >= '0' && r <= '9' {
						m.input += string(r)
					}
				}
			}
		case data.StatePlayerAction:
			if text == "?" {
				m.showHelp()
				m.err = nil
				break
			}
			var command string
			switch {
			case text == "h":
				command = "hit"
			case text == "s" || key.Code == tea.KeyEnter:
				command = "stand"
			case text == "d":
				command = "double"
			case text == "p":
				command = "split"
			default:
				return m, nil
			}
			if err := m.handleCommand(command); err != nil {
				m.err = err
			} else {
				m.err = nil
			}
		default:
			// ignore input in other states
		}
	}
	return m, nil
}

func (m *Model) View() string {
	if m.quitting {
		return "Thanks for playing!\n"
	}

	header := headerStyle.Render("♣ Blackjack")
	info := infoStyle.Render(fmt.Sprintf("Deck cards remaining: %d", m.game.Deck().CardsLeft()))

	dealerSection := m.renderDealerSection()
	playerSection := m.renderPlayerSection()
	hotkeys := m.renderHotkeys()
	prompt := m.renderPromptArea()

	sections := []string{header, info, dealerSection, playerSection}
	if hotkeys != "" {
		sections = append(sections, hotkeys)
	}
	if prompt != "" {
		sections = append(sections, prompt)
	}
	if len(m.messages) > 0 {
		sections = append(sections, m.renderMessages())
	}
	if m.err != nil {
		sections = append(sections, errorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m *Model) handleCommand(cmd string) error {
	switch m.game.State() {
	case data.StateBetting, data.StateSettled:
		if m.game.State() == data.StateSettled {
			m.game.PrepareNextRound()
		}
		amount, err := strconv.Atoi(cmd)
		if err != nil {
			return fmt.Errorf("invalid bet amount: %w", err)
		}
		if amount <= 0 {
			return fmt.Errorf("bet must be positive")
		}
		bets := map[string]int{}
		if m.player != nil {
			bets[m.player.Name()] = amount
		}
		if err := m.game.StartRound(bets); err != nil {
			return err
		}
		m.results = nil
		m.messages = nil
		m.log(fmt.Sprintf("Bet $%d", amount))
		if err := m.game.DealInitialCards(); err != nil {
			return err
		}
		m.log("Cards dealt")
		if m.player != nil {
			hand := m.player.ActiveHand()
			if hand.IsBlackjack() {
				m.log("Blackjack!")
				hand.Stand()
				m.player.SetStatus(data.PlayerStatusStanding)
				if m.game.ReadyForDealer() {
					return m.completeRound()
				}
			}
		}
		m.updatePrompt()
		return nil
	case data.StatePlayerAction:
		if m.player == nil {
			return fmt.Errorf("no player available")
		}
		hand := m.player.ActiveHand()
		lower := strings.ToLower(cmd)
		switch lower {
		case "hit":
			if hand == nil {
				return data.ErrNoActiveHand
			}
			if hand.IsStanding() {
				return fmt.Errorf("hand already standing")
			}
			card, err := m.game.Hit(m.player)
			if err != nil {
				return err
			}
			m.log(fmt.Sprintf("Hit: drew %s", card.String()))
			if hand.IsBusted() {
				m.log(fmt.Sprintf("Busted with %d", hand.Value()))
			}
			if m.game.ReadyForDealer() {
				return m.completeRound()
			}
		case "stand":
			if hand == nil {
				return data.ErrNoActiveHand
			}
			if hand.IsStanding() {
				return fmt.Errorf("hand already standing")
			}
			if err := m.game.Stand(m.player); err != nil {
				return err
			}
			m.log("Stand")
			if m.game.ReadyForDealer() {
				return m.completeRound()
			}
		case "double":
			if hand == nil {
				return data.ErrNoActiveHand
			}
			if hand.IsStanding() {
				return fmt.Errorf("hand already standing")
			}
			if err := m.player.DoubleDownActiveHand(); err != nil {
				return err
			}
			card, err := m.game.Hit(m.player)
			if err != nil {
				return err
			}
			m.log(fmt.Sprintf("Double down: drew %s", card.String()))
			if hand.IsBusted() {
				m.log(fmt.Sprintf("Busted with %d", hand.Value()))
			}
			if err := m.game.Stand(m.player); err != nil {
				return err
			}
			if m.game.ReadyForDealer() {
				return m.completeRound()
			}
		case "split":
			if hand == nil {
				return data.ErrNoActiveHand
			}
			newHand, err := m.player.SplitActiveHand()
			if err != nil {
				return err
			}
			firstCard := m.game.Deck().Deal()
			hand.AddCard(firstCard)
			secondCard := m.game.Deck().Deal()
			newHand.AddCard(secondCard)
			m.log(fmt.Sprintf("Split hand. Drew %s and %s", firstCard.String(), secondCard.String()))
		default:
			return fmt.Errorf("unknown command: %s", cmd)
		}
		m.updatePrompt()
		return nil
	default:
		return fmt.Errorf("game not ready for input")
	}
}

func (m *Model) completeRound() error {
	if err := m.game.DealerPlay(); err != nil {
		return err
	}
	results, err := m.game.SettleRound()
	if err != nil {
		return err
	}
	m.results = results
	m.messages = nil
	for _, res := range results {
		m.log(fmt.Sprintf("%s %s", res.Player.Name(), describeOutcome(res)))
	}
	m.updatePrompt()
	return nil
}

func (m *Model) renderDealerSection() string {
	dealer := m.game.Dealer()
	hand := dealer.ActiveHand()
	title := sectionTitleStyle.Render("Dealer")

	var cards []string
	if hand != nil && len(hand.Cards()) > 0 {
		if dealer.HoleCardHidden() {
			cards = append(cards, renderCard(hand.Cards()[0]))
			cards = append(cards, renderFacedownCard())
		} else {
			for _, card := range hand.Cards() {
				cards = append(cards, renderCard(card))
			}
		}
	}

	var body string
	if len(cards) == 0 {
		body = infoStyle.Render("Waiting to deal…")
	} else {
		body = lipgloss.JoinHorizontal(lipgloss.Top, cards...)
	}

	valueText := ""
	if hand != nil {
		if dealer.HoleCardHidden() {
			valueText = infoStyle.Render("Value: ??")
		} else {
			valueText = valueStyle.Render(fmt.Sprintf("Value: %d", hand.Value()))
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, title, body, valueText)
}

func (m *Model) renderPlayerSection() string {
	if m.player == nil {
		return ""
	}

	header := sectionTitleStyle.Render(fmt.Sprintf("%s — Bankroll: $%d", m.player.Name(), m.player.Bankroll()))
	var handViews []string
	for i, hand := range m.player.Hands() {
		active := m.game.State() == data.StatePlayerAction && i == m.player.ActiveHandIndex()
		handViews = append(handViews, renderPlayerHand(hand, i, active))
	}
	if len(handViews) == 0 {
		handViews = append(handViews, infoStyle.Render("No cards yet"))
	}
	return lipgloss.JoinVertical(lipgloss.Left, header, lipgloss.JoinVertical(lipgloss.Left, handViews...))
}

func (m *Model) renderHotkeys() string {
	switch m.game.State() {
	case data.StatePlayerAction:
		hand := m.player.ActiveHand()
		hotkeys := []hotkey{
			{Key: "H", Label: "Hit", Enabled: hand != nil && !hand.IsStanding() && !hand.IsBusted()},
			{Key: "S", Label: "Stand", Enabled: hand != nil && !hand.IsStanding()},
			{Key: "D", Label: "Double", Enabled: canDouble(m.player, hand)},
			{Key: "P", Label: "Split", Enabled: canSplit(m.player, hand)},
			{Key: "?", Label: "Help", Enabled: true},
			{Key: "Q", Label: "Quit", Enabled: true},
		}
		return hotkeyBarStyle.Render(renderHotkeyLine(hotkeys))
	case data.StateBetting, data.StateSettled:
		hotkeys := []hotkey{
			{Key: "?", Label: "Help", Enabled: true},
			{Key: "Q", Label: "Quit", Enabled: true},
		}
		return hotkeyBarStyle.Render(renderHotkeyLine(hotkeys))
	default:
		return ""
	}
}

func (m *Model) renderPromptArea() string {
	switch m.game.State() {
	case data.StateBetting:
		return lipgloss.JoinVertical(lipgloss.Left,
			promptStyle.Render("Bet amount (press Enter to confirm):"),
			inputStyle.Render(fmt.Sprintf("$%s", m.input)))
	case data.StateSettled:
		return lipgloss.JoinVertical(lipgloss.Left,
			promptStyle.Render("Round settled. Enter next bet or press Q to quit."),
			inputStyle.Render(fmt.Sprintf("$%s", m.input)))
	case data.StatePlayerAction:
		return promptStyle.Render("Hotkeys: [H]it [S]tand [D]ouble [P]Split [?]Help [Q]Quit")
	default:
		return ""
	}
}

func (m *Model) renderMessages() string {
	var lines []string
	if len(m.results) > 0 {
		lines = append(lines, "Last round:")
		for _, res := range m.results {
			lines = append(lines, fmt.Sprintf("  %s", describeOutcome(res)))
		}
	}
	if len(m.messages) > 0 {
		lines = append(lines, "Messages:")
		for _, message := range m.messages {
			lines = append(lines, fmt.Sprintf("  %s", message))
		}
	}
	return messageBoxStyle.Render(strings.Join(lines, "\n"))
}

func (m *Model) updatePrompt() {
	switch m.game.State() {
	case data.StateBetting:
		m.prompt = "Enter bet amount"
	case data.StatePlayerAction:
		m.prompt = "Hotkeys: [H]it [S]tand [D]ouble [P]Split"
	case data.StateSettled:
		m.prompt = "Round settled"
	default:
		m.prompt = ""
	}
}

func (m *Model) log(message string) {
	if message == "" {
		return
	}
	m.messages = append(m.messages, message)
	const maxMessages = 6
	if len(m.messages) > maxMessages {
		m.messages = m.messages[len(m.messages)-maxMessages:]
	}
}

func (m *Model) showHelp() {
	help := []string{
		"Bet: type numbers then press Enter.",
		"Hotkeys during play: H=Hit, S=Stand, D=Double, P=Split.",
		"Press ? anytime to show this help, Q to quit, Ctrl+C also exits.",
	}
	for _, line := range help {
		m.log(line)
	}
}

type hotkey struct {
	Key     string
	Label   string
	Enabled bool
}

func renderHotkeyLine(hotkeys []hotkey) string {
	var parts []string
	for _, hk := range hotkeys {
		keyStyle := hotkeyKeyStyle
		if !hk.Enabled {
			keyStyle = hotkeyDisabledKey
		}
		parts = append(parts, keyStyle.Render("["+hk.Key+"] ")+hotkeyLabelStyle.Render(hk.Label))
	}
	return strings.Join(parts, "  ")
}

func renderPlayerHand(hand *data.Hand, index int, active bool) string {
	if hand == nil {
		return ""
	}
	title := fmt.Sprintf("Hand %d", index+1)
	cards := hand.Cards()
	var rendered []string
	for _, card := range cards {
		rendered = append(rendered, renderCard(card))
	}
	if len(rendered) == 0 {
		rendered = append(rendered, infoStyle.Render("(empty)"))
	}
	cardRow := lipgloss.JoinHorizontal(lipgloss.Top, rendered...)

	var tags []string
	if hand.IsBlackjack() {
		tags = append(tags, "BLACKJACK")
	}
	if hand.IsSoft() && !hand.IsBlackjack() {
		tags = append(tags, "SOFT")
	}
	if hand.IsBusted() {
		tags = append(tags, "BUST")
	}
	if hand.IsDoubleDown() {
		tags = append(tags, "DOUBLE")
	}
	if hand.IsStanding() {
		tags = append(tags, "STAND")
	}

	tagText := ""
	if len(tags) > 0 {
		tagText = tagStyle.Render(strings.Join(tags, " · "))
	}

	info := fmt.Sprintf("Value: %d", hand.Value())
	if bet := hand.Bet(); bet > 0 {
		info += fmt.Sprintf("   Bet: $%d", bet)
	}

	box := lipgloss.JoinVertical(lipgloss.Left,
		sectionTitleStyle.Render(title),
		cardRow,
		valueStyle.Render(info),
		tagText,
	)

	if active {
		return activeHandStyle.Render(box)
	}
	return handBoxStyle.Render(box)
}

func renderCard(card data.Card) string {
	rank := data.RankString[card.Rank]
	suit := data.SuitString[card.Suit]
	color := suitColor(card.Suit)
	lines := []string{
		fmt.Sprintf("%-2s      ", rank),
		"         ",
		fmt.Sprintf("   %s   ", suit),
		"         ",
		fmt.Sprintf("      %-2s", rank),
	}
	content := strings.Join(lines, "\n")
	return cardStyle.Copy().Foreground(color).BorderForeground(color).Render(content)
}

func renderFacedownCard() string {
	shade := strings.Repeat("▓", 7)
	lines := []string{shade, shade, shade, shade, shade}
	return faceDownStyle.Render(strings.Join(lines, "\n"))
}

func suitColor(suit data.Suit) color.Color {
	switch suit {
	case data.Hearts:
		return data.HeartColor
	case data.Diamonds:
		return data.DiamondColor
	case data.Clubs:
		return data.ClubColor
	case data.Spades:
		fallthrough
	default:
		return data.SpadeColor
	}
}

func canDouble(player *data.Player, hand *data.Hand) bool {
	if player == nil || hand == nil {
		return false
	}
	if hand.IsStanding() || hand.IsBusted() || hand.IsDoubleDown() {
		return false
	}
	if len(hand.Cards()) != 2 {
		return false
	}
	bet := hand.Bet()
	if bet == 0 {
		return false
	}
	return player.Bankroll() >= bet
}

func canSplit(player *data.Player, hand *data.Hand) bool {
	if player == nil || hand == nil {
		return false
	}
	if hand.IsStanding() || hand.IsBusted() {
		return false
	}
	if !hand.CanSplit() {
		return false
	}
	return player.Bankroll() >= hand.Bet()
}

func describeOutcome(res data.RoundResult) string {
	hand := res.Hand
	value := 0
	if hand != nil {
		value = hand.Value()
	}
	switch res.Outcome {
	case data.OutcomeWin:
		return fmt.Sprintf("wins with %d", value)
	case data.OutcomeBlackjack:
		return "wins with blackjack"
	case data.OutcomePush:
		return fmt.Sprintf("push with %d", value)
	default:
		return fmt.Sprintf("loses with %d", value)
	}
}

func trimLastRune(s string) string {
	if s == "" {
		return s
	}
	_, size := utf8.DecodeLastRuneInString(s)
	return s[:len(s)-size]
}
