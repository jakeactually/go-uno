package models

import "math/rand"

type CardColor string

const (
	Red    CardColor = "red"
	Green            = "green"
	Blue             = "blue"
	Yellow           = "yellow"
)

type CardType int

const (
	Number      CardType = 0
	Stop                 = 1
	Reverse              = 2
	Plus2                = 3
	ChangeColor          = 4
	Plus4                = 5
)

type Card = struct {
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

		cards = append(cards, Card{GetColor: c, GetType: Stop})
		cards = append(cards, Card{GetColor: c, GetType: Reverse})
		cards = append(cards, Card{GetColor: c, GetType: Plus2})
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
