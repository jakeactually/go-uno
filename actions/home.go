package actions

import (
	"go_uno/models"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("index.html"))
}

// NewRoomHandler
func NewRoomHandler(c buffalo.Context) error {
	conn, _ := pop.Connect("development")
	room := &models.Room{}
	conn.Create(room)

	return c.Render(http.StatusOK, r.HTML("index.html"))
}
