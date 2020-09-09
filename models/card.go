package models

import (
	"fmt"
	"math/rand"

	"github.com/gobuffalo/pop/nulls"
)

type CardColor string

const (
	Red    CardColor = "red"
	Green            = "green"
	Blue             = "blue"
	Yellow           = "yellow"
)

type CardType string

const (
	Number      CardType = "number"
	Stop                 = "stop"
	Reverse              = "reverse"
	Plus2                = "plus2"
	ChangeColor          = "color"
	Plus4                = "plus4"
)

type Card struct {
	GetNumber int
	GetColor  CardColor
	GetType   CardType
}

var Colors = []CardColor{Red, Green, Blue, Yellow}

func MakeAllCards() []Card {
	var cards []Card

	for n := range [9]int{} {
		for _, c := range Colors {
			for range [2]int{} {
				cards = append(cards, Card{
					GetNumber: n + 1,
					GetColor:  c,
					GetType:   Number,
				})
			}
		}
	}

	for _, c := range Colors {
		// There are 4 zeroes
		cards = append(cards, Card{GetColor: c, GetType: Number})

		for range [2]int{} {
			cards = append(cards, Card{GetColor: c, GetType: Stop})
			cards = append(cards, Card{GetColor: c, GetType: Reverse})
			cards = append(cards, Card{GetColor: c, GetType: Plus2})
		}
	}

	for range [4]int{} {
		// The color of these cards is ignored
		cards = append(cards, Card{GetColor: Red, GetType: ChangeColor})
		cards = append(cards, Card{GetColor: Red, GetType: Plus4})
	}

	return cards
}

var AllCards = MakeAllCards()

func Deck() []int {
	var deck []int

	for i := range [108]int{} {
		deck = append(deck, i)
	}

	return deck
}

func Shuffle(deck []int) {
	for range deck {
		a := rand.Intn(len(deck))
		b := rand.Intn(len(deck))

		temp := deck[a]
		deck[a] = deck[b]
		deck[b] = temp
	}
}

func (c1 Card) Matches(gameState nulls.String, chosenColor CardColor, c2 Card) (bool, string) {
	if gameState.Valid {
		if gameState.String == "stop" {
			return c2.GetType == Stop, "You can only chain or pass"
		} else if gameState.String == "plus2" {
			return c2.GetType == Plus2, "You can only chain or draw"
		} else if gameState.String == "plus4" {
			return c2.GetType == Plus4, "You can only chain or draw"
		}
	}

	return c1.FreeMatch(chosenColor, c2)
}

func (c1 Card) FreeMatch(chosenColor CardColor, c2 Card) (bool, string) {
	if c1.GetType == Number {
		if c2.GetType == Number {
			return c1.GetNumber == c2.GetNumber || c1.GetColor == c2.GetColor, "Invalid Move"
		} else if c2.IsSign() {
			return c1.GetColor == c2.GetColor, "Wrong color"
		}

		return true, ""
	} else if c1.IsSign() {
		if c2.GetType == Number {
			return c1.GetColor == c2.GetColor, "Wrong color"
		} else if c2.IsSign() {
			return c1.GetType == c2.GetType || c1.GetColor == c2.GetColor, "Invalid Move"
		}

		return true, ""
	}

	if c2.GetType == Number || c2.IsSign() {
		return chosenColor == c2.GetColor, "Choosen color is " + string(chosenColor)
	}

	return true, ""
}

func (c Card) IsSign() bool {
	return c.GetType == Stop || c.GetType == Reverse || c.GetType == Plus2
}

func (c Card) CanChain() bool {
	return c.GetType == Stop || c.GetType == Plus2 || c.GetType == Plus4
}

func (c Card) IsColorCard() bool {
	return c.GetType == ChangeColor || c.GetType == Plus4
}

func (c Card) Image() string {
	if c.GetType == Number {
		return fmt.Sprintf("%c%d", c.GetColor[0], c.GetNumber)
	} else if c.IsSign() {
		return fmt.Sprintf("%c-%s", c.GetColor[0], c.GetType)
	}

	return string(c.GetType)
}
