package actions

import (
	"go_uno/models"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/pop"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("index.html"))
}

// NewRoomHandler ...
func NewRoomHandler(c buffalo.Context) error {
	conn, _ := pop.Connect("development")

	room := &models.Room{}
	conn.Create(room)

	player := &models.Player{Name: c.Param("username"), RoomID: room.ID}
	conn.Create(player)
	c.Request().Method = "GET"

	return c.Redirect(307, "roomPath()", render.Data{"roomID": room.ID})
}

// RoomHandler ...
func RoomHandler(c buffalo.Context) error {
	conn, _ := pop.Connect("development")

	room := &models.Room{}
	conn.Find(room, c.Param("roomID"))
	c.Set("room", room)

	return c.Render(http.StatusOK, r.HTML("room.plush.html"))
}
