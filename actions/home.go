package actions

import (
	"go_uno/models"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/pop"
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
	conn.Load(room)
	c.Set("room", room)

	return c.Render(http.StatusOK, r.HTML("room.plush.html"))
}

// RoomStateHandler ...
func RoomStateHandler(c buffalo.Context) error {
	ws, err := websocket.Upgrade(c.Response(), c.Request(), c.Response().Header(), 1024, 1024)

	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			errors.WithStack(err)
		}
	}

	client := &Client{hub: hub, conn: ws, uid: uid, name: name, buffer: make(chan []byte, 256)}
	client.hub.register <- client

	go client.Reader()
	client.Writer()

	client.hub.unregister <- client
	client.conn.Close()

	return nil
}
