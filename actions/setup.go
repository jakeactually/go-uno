package actions

import (
	"go_uno/models"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
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
			hand, cards = cards[0:7], cards[7:]
			p.Cards = hand
			conn.Update(&p)
		}

		room.Deck = cards
		// hack
		// buffalo pop loads postgres empty int arrays as slices with a single zero
		room.Center = []int{}
		room.ToCenter()

		for models.AllCards[room.Top()].GetType != models.Number {
			room.ToCenter()
		}

		room.Active = true
		conn.Update(room)
	}

	c.Set("room", room)

	return c.Render(http.StatusOK, r.HTML("game.plush.html"))
}

// CenterHandler ...
func CenterHandler(c buffalo.Context) error {
	conn, _ := pop.Connect("development")

	room := &models.Room{}
	conn.Find(room, c.Param("roomID"))

	cards := []interface{}{}
	for _, c := range room.Center {
		cards = append(cards, []interface{}{c, models.AllCards[c].Image()})
	}

	return c.Render(http.StatusOK, r.JSON(cards))
}

// HandHandler ...
func HandHandler(c buffalo.Context) error {
	conn, _ := pop.Connect("development")

	s := c.Session()
	pid, _ := s.Get("playerID").(uuid.UUID)
	player := &models.Player{}
	conn.Find(player, pid)

	cards := []interface{}{}
	for _, c := range noZeroes(player.Cards) {
		cards = append(cards, []interface{}{c, models.AllCards[c].Image()})
	}

	return c.Render(http.StatusOK, r.JSON(cards))
}

// GameOverHandler ...
func GameOverHandler(c buffalo.Context) error {
	conn, _ := pop.Connect("development")

	room := &models.Room{}
	conn.Find(room, c.Param("roomID"))
	room.Active = false
	conn.Update(room)

	c.Set("room", room)

	return c.Render(http.StatusOK, r.HTML("theEnd.plush.html"))
}
