package actions

import (
	"fmt"
	"go_uno/models"
	"log"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	s := c.Session()
	playerID := s.Get("playerID")
	conn, _ := pop.Connect("development")
	player := &models.Player{}

	if playerID != nil {
		pid, _ := playerID.(uuid.UUID)
		conn.Find(player, pid)
	}

	c.Set("player", player)

	return c.Render(http.StatusOK, r.HTML("index.html"))
}

// NewRoomHandler ...
func NewRoomHandler(c buffalo.Context) error {
	conn, _ := pop.Connect("development")

	room := &models.Room{}
	conn.Create(room)
	player := &models.Player{Name: c.Param("username"), RoomID: room.ID}
	conn.Create(player)

	s := c.Session()
	s.Set("playerID", player.ID)
	s.Save()

	return c.Redirect(http.StatusMovedPermanently, "roomPath()", render.Data{"roomID": room.ID})
}

// RoomHandler ...
func RoomHandler(c buffalo.Context) error {
	conn, _ := pop.Connect("development")

	room := &models.Room{}
	conn.Find(room, c.Param("roomID"))

	if room.Active {
		return c.Render(http.StatusBadRequest, r.String("This room is already playing"))
	}

	conn.Load(room)

	s := c.Session()
	playerID := s.Get("playerID")
	pid, _ := playerID.(uuid.UUID)
	isInRoom := false

	for _, p := range room.Players {
		if p.ID == pid {
			isInRoom = true
		}
	}

	if !isInRoom {
		return c.Redirect(http.StatusMovedPermanently, "joinRoomPath()", render.Data{"roomID": room.ID})
	}

	player := &models.Player{}
	conn.Find(player, pid)

	c.Set("room", room)
	c.Set("player", player)

	return c.Render(http.StatusOK, r.HTML("room.plush.html"))
}

// JoinRoomHandler ...
func JoinRoomHandler(c buffalo.Context) error {
	s := c.Session()
	playerID := s.Get("playerID")
	conn, _ := pop.Connect("development")
	player := &models.Player{}

	if playerID != nil {
		pid, _ := playerID.(uuid.UUID)
		conn.Find(player, pid)
	}

	c.Set("roomId", c.Param("roomID"))
	c.Set("player", player)

	return c.Render(http.StatusOK, r.HTML("join.plush.html"))
}

// JoinRoomPostHandler ...
func JoinRoomPostHandler(c buffalo.Context) error {
	conn, _ := pop.Connect("development")

	room := &models.Room{}
	conn.Find(room, c.Param("roomID"))
	player := &models.Player{Name: c.Param("username"), RoomID: room.ID}
	conn.Create(player)

	s := c.Session()
	s.Set("playerID", player.ID)
	s.Save()

	for _, ws := range notifier[fmt.Sprintf("%d", room.ID)] {
		ws.WriteMessage(websocket.TextMessage, []byte("update"))
	}

	return c.Redirect(http.StatusMovedPermanently, "roomPath()", render.Data{"roomID": room.ID})
}

// ExpelHandler ...
func ExpelHandler(c buffalo.Context) error {
	conn, _ := pop.Connect("development")

	room := &models.Room{}
	conn.Find(room, c.Param("roomID"))
	conn.Load(room)

	for _, p := range room.Players {
		log.Println(p.ID.String(), " ", c.Param("playerID"))
		if p.ID.String() == c.Param("playerID") {
			p.RoomID = 0
			conn.Update(&p)
		}
	}

	for _, ws := range notifier[fmt.Sprintf("%d", room.ID)] {
		ws.WriteMessage(websocket.TextMessage, []byte("update"))
	}

	return c.Redirect(http.StatusMovedPermanently, "rootPath()")
}

var notifier = make(map[string][](*websocket.Conn))
var upgrader = websocket.Upgrader{}

// RoomStateHandler ...
func RoomStateHandler(c buffalo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)

	if err != nil {
		return err
	}

	notifier[c.Param("roomID")] = append(notifier[c.Param("roomID")], ws)

	return nil
}
