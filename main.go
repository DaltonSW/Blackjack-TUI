package main

import (
	"blackjack/internal/data"
	"fmt"
)

func main() {
	style := data.CardStyle.Foreground(data.SpadeColor)
	fmt.Println(style.Render(fmt.Sprintf("%s\n%s\n%s", "5  ", " "+data.SuitString[data.Spades]+" ", "  5")))
}
