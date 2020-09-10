package actions

import (
	"fmt"
	"go_uno/models"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
)

var gameNotifier = make(map[string][](*websocket.Conn))
var gameUpgrader = websocket.Upgrader{}

// TurnHandler ...
func TurnHandler(c buffalo.Context) error {
	conn, _ := pop.Connect("development")
	room, player := roomPlayer(c, conn)

	if room.Players[room.Turn].ID != player.ID {
		return c.Render(http.StatusBadRequest, r.String("Not your turn"))
	}

	cardID, _ := strconv.Atoi(c.Param("cardId"))
	topCardID := room.Top()
	card1 := models.AllCards[topCardID]
	card2 := models.AllCards[cardID]

	chosenColor := models.Red
	if card1.IsColorCard() {
		chosenColor = models.CardColor(room.Color)
	}

	ok, err := card1.Matches(room.GameState, chosenColor, card2)

	if ok {
		effects(c, room, card2)

		player.Cards = intFilter(player.Cards, cardID)
		room.Center = append(room.Center, cardID)

		room.Next()
		player.Drawed = false

		conn.Update(room)
		conn.Update(player)

		// Notify
		for _, ws := range gameNotifier[strconv.Itoa(room.ID)] {
			ws.WriteMessage(websocket.TextMessage, []byte("update"))
		}

		return c.Render(http.StatusOK, nil)
	}

	return c.Render(http.StatusBadRequest, r.String(err))
}

func effects(c buffalo.Context, room *models.Room, card models.Card) {
	if card.IsColorCard() {
		room.Color = c.Param("color")
	}

	if card.CanChain() {
		room.GameState = nulls.String{String: string(card.GetType), Valid: true}
		room.ChainCount++
	} else {
		room.ChainCount = 0
	}

	if card.GetType == models.Reverse {
		room.Direction = !room.Direction
	}
}

// GameHandler ...
func GameHandler(c buffalo.Context) error {
	ws, err := gameUpgrader.Upgrade(c.Response(), c.Request(), nil)

	if err != nil {
		return err
	}

	gameNotifier[c.Param("roomID")] = append(gameNotifier[c.Param("roomID")], ws)

	return nil
}

// DrawHandler ...
func DrawHandler(c buffalo.Context) error {
	conn, _ := pop.Connect("development")
	room, player := roomPlayer(c, conn)

	player.Draw(room)
	player.Drawed = true

	conn.Update(room)
	conn.Update(player)

	return c.Render(http.StatusOK, nil)
}

func roomPlayer(c buffalo.Context, conn *pop.Connection) (*models.Room, *models.Player) {
	room := &models.Room{}
	conn.Find(room, c.Param("roomID"))
	conn.Load(room)

	s := c.Session()
	pid, _ := s.Get("playerID").(uuid.UUID)
	player := &models.Player{}
	conn.Find(player, pid)

	return room, player
}

// BoardStateHandler ...
func BoardStateHandler(c buffalo.Context) error {
	conn, _ := pop.Connect("development")
	room, player := roomPlayer(c, conn)

	isTurn := room.Players[room.Turn].ID == player.ID
	obj := []interface{}{isTurn, player.Drawed, room.Color, room.GameState, room.ChainCount}

	return c.Render(http.StatusOK, r.JSON(obj))
}

// PassHandler ...
func PassHandler(c buffalo.Context) error {
	conn, _ := pop.Connect("development")
	room, player := roomPlayer(c, conn)

	if player.Drawed {
		room.Next()
		player.Drawed = false
		conn.Update(room)
		conn.Update(player)

		// Notify
		for _, ws := range gameNotifier[fmt.Sprintf("%d", room.ID)] {
			ws.WriteMessage(websocket.TextMessage, []byte("update"))
		}

		return c.Render(http.StatusOK, nil)
	}

	return c.Render(http.StatusBadRequest, r.String("You must draw one card"))
}

// PenaltyHandler ...
func PenaltyHandler(c buffalo.Context) error {
	conn, _ := pop.Connect("development")
	room, player := roomPlayer(c, conn)

	if room.GameState.String == "plus2" {
		penalty(room, player, 2*room.ChainCount)
	} else if room.GameState.String == "plus4" {
		penalty(room, player, 4*room.ChainCount)
	} else if room.GameState.String == "stop" {
		room.GameState = nulls.String{String: "", Valid: false}
		room.Next()
	}

	conn.Update(room)
	conn.Update(player)

	// Notify
	for _, ws := range gameNotifier[fmt.Sprintf("%d", room.ID)] {
		ws.WriteMessage(websocket.TextMessage, []byte("update"))
	}

	return c.Render(http.StatusOK, nil)
}

func penalty(room *models.Room, player *models.Player, amount int) {
	for i := 0; i < amount; i++ {
		player.Draw(room)
	}

	room.GameState = nulls.String{String: "", Valid: false}
	room.ChainCount = 0
	room.Next()
}

// AllPlayersHandler ...
func AllPlayersHandler(c buffalo.Context) error {
	conn, _ := pop.Connect("development")
	room := &models.Room{}
	conn.Find(room, c.Param("roomID"))
	conn.Load(room)

	var obj []interface{}

	for i, p := range room.Players {
		obj = append(obj, []interface{}{twoLetters(p.Name), len(noZeroes(p.Cards)), room.Turn == i})
	}

	return c.Render(http.StatusOK, r.JSON(obj))
}

// hack
// buffalo pop loads postgres empty int arrays as slices with a single zero
func noZeroes(arr []int) []int {
	return intFilter(arr, 0)
}

func intFilter(arr []int, n int) []int {
	out := []int{}

	for _, x := range arr {
		if x == n {
			continue
		}

		out = append(out, x)
	}

	return out
}

func twoLetters(fullName string) string {
	regex, _ := regexp.Compile("\\b\\w")
	initials := regex.FindAllString(fullName, -1)

	var letters []string

	if len(initials) >= 2 {
		letters = initials
	} else {
		regex, _ := regexp.Compile("\\w")
		letters = regex.FindAllString(fullName, -1)
	}

	l := len(letters)
	var rl int

	if l > 2 {
		rl = 2
	} else {
		rl = l
	}

	return strings.Join(letters[0:rl], "")
}
