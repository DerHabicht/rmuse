package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/satori/go.uuid"
)

type Medium struct {
	ID         uuid.UUID `json:"id" db:"id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	URI        string    `json:"uri"        db:"uri"`
	User       uuid.UUID `json:"userid"     db:"user_id"`
	Filetype   string    `json:"type"       db:"filetype"`
	Permission string    `json:"permission" db:"permission"`
}

func (m *Medium) Create(tx *pop.Connection) (*validate.Errors, error) {
	return tx.ValidateAndCreate(m)
}

func GetMediumByID(tx *pop.Connection, id uuid.UUID) (*Medium, error) {
	m := Medium{}
	err := tx.Find(&m, id)

	if err != nil {
		return nil, fmt.Errorf("could not find media %v", err)
	}

	return &m, nil
}

// String is not required by pop and may be deleted
func (m Medium) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

// Media is not required by pop and may be deleted
type Media []Medium

// String is not required by pop and may be deleted
func (m Media) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (m *Medium) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (m *Medium) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	var err error
	return validate.Validate(
		&validators.FuncValidator{
			Field:   m.Filetype,
			Name:    "Filetype",
			Message: "type is empty",
			Fn: func() bool {
				return m.Filetype != ""
			},
		},
		&validators.FuncValidator{
			Field:   m.URI,
			Name:    "URI",
			Message: "there is already a file with URI %s",
			Fn: func() bool {
				var b bool
				q := tx.Where("uri = ?", m.URI)
				b, err = q.Exists(m)
				if err != nil {
					return false
				}
				return !b
			},
		},
		&validators.FuncValidator{
			Field:   m.URI,
			Name:    "URI",
			Message: "uri is empty",
			Fn: func() bool {
				return m.URI != ""
			},
		},
	), err
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (m *Medium) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
