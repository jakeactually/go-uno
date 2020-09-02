package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/slices"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

// Player is used by pop to map your players database table to your go code.
type Player struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	Name   string     `json:"name" db:"name"`
	Cards  slices.Int `json:"cards" db:"cards"`
	Drawed bool       `json:"drawed" db:"drawed"`

	RoomID int `json:"room_id" db:"room_id"`
}

// String is not required by pop and may be deleted
func (p Player) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Players is not required by pop and may be deleted
type Players []Player

// String is not required by pop and may be deleted
func (p Players) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (p *Player) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (p *Player) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (p *Player) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
