package main

import "fmt"

func main() {
	style := CardStyle.Foreground(SpadeColor)
	fmt.Println(style.Render(fmt.Sprintf("%s\n%s\n%s", "5  ", " "+SpadeSymbol+" ", "  5")))
}
