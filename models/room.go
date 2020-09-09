package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/nulls"
	"github.com/gobuffalo/pop/slices"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
)

// Room is used by pop to map your rooms database table to your go code.
type Room struct {
	ID        int       `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	Active     bool         `json:"active" db:"active"`
	Deck       slices.Int   `json:"deck" db:"deck"`
	Center     slices.Int   `json:"center" db:"center"`
	Turn       int          `json:"turn" db:"turn"`
	Direction  bool         `json:"direction" db:"direction"`
	Color      string       `json:"color" db:"color"`
	GameState  nulls.String `json:"gameState" db:"game_state"`
	ChainCount int          `json:"chainCount" db:"chain_count"`

	Players []Player `json:"players" has_many:"players" order_by:"ID"`
}

// String is not required by pop and may be deleted
func (r Room) String() string {
	jr, _ := json.Marshal(r)
	return string(jr)
}

// Rooms is not required by pop and may be deleted
type Rooms []Room

// String is not required by pop and may be deleted
func (r Rooms) String() string {
	jr, _ := json.Marshal(r)
	return string(jr)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (r *Room) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (r *Room) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (r *Room) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (room *Room) ToCenter() {
	room.Center, room.Deck = append(room.Center, room.Deck[0]), room.Deck[1:]
}

func (room *Room) DrawOne() int {
	result := room.Deck[0]
	room.Deck = room.Deck[1:]
	return result
}

func (room Room) Top() int {
	return room.Center[len(room.Center)-1]
}

func (room *Room) Left() {
	if room.Turn == 0 {
		room.Turn = len(room.Players) - 1
	} else {
		room.Turn--
	}
}

func (room *Room) Right() {
	if room.Turn == len(room.Players)-1 {
		room.Turn = 0
	} else {
		room.Turn++
	}
}

func (room *Room) Next() {
	if room.Direction {
		room.Right()
	} else {
		room.Left()
	}

	player := room.Players[room.Turn]

	if len(player.Cards) == 0 {
		room.Next()
	}
}
