package actions

import (
	"go_uno/models"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
)

// PlayHandler ...
func PlayHandler(c buffalo.Context) error {
	conn, _ := pop.Connect("development")

	room := &models.Room{}
	conn.Find(room, c.Param("roomID"))
	conn.Load(room)

	if len(room.Players) < 2 {
		return c.Render(http.StatusBadRequest, r.String("Not enough players"))
	}

	if !room.Active {
		var hand []int
		cards := models.Deck()
		models.Shuffle(cards)

		for _, p := range room.Players {
			hand, cards = cards[0:6], cards[7:]
			p.Cards = hand
			conn.Update(&p)
		}

		room.Deck = cards
		room.Center, room.Deck = append(room.Center, room.Deck[0]), room.Deck[1:]

		for models.AllCards[room.Center[len(room.Center)-1]].GetType != models.Number {
			room.Center, room.Deck = append(room.Center, room.Deck[0]), room.Deck[1:]
		}

		room.Active = true
		conn.Update(room)
	}

	c.Set("room", room)

	return c.Render(http.StatusOK, r.HTML("game.plush.html"))
}
