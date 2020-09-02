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

	Players []Player `json:"players" has_many:"players"`
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
