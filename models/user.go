package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/envy"
	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/satori/go.uuid"
)

type User struct {
	ID           uuid.UUID `json:"user_id"    db:"id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	Email        string    `json:"email"      db:"email"`
	Username     string    `json:"username"   db:"username"`
	FirstName    string    `json:"firstname"  db:"first_name"`
	LastName     string    `json:"lastname"   db:"last_name"`
	UserType     string    `json:"type"       db:"role"`
	PasswordHash string    `json:"-"          db:"password_hash"`
}

func (u *User) CreateJWTToken() (string, error) {
	// Create and return a JWT token
	exp, _ := time.ParseDuration("168h")
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(exp).Unix(),
		Id:        u.ID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signingKey, err := ioutil.ReadFile(envy.Get("JWT_KEY_PATH", "jwtRS256.key"))
	if err != nil {
		return "", fmt.Errorf("could not open jwt key, %v", err)
	}

	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", fmt.Errorf("could sign token, %v", err)
	}

	return tokenString, nil
}

func (u *User) Create(tx *pop.Connection) (*validate.Errors, error) {
	u.Email = strings.ToLower(u.Email)

	return tx.ValidateAndCreate(u)
}

func (u *User) Update(tx *pop.Connection) (*validate.Errors, error) {
	u.Email = strings.ToLower(u.Email)

	return tx.ValidateAndUpdate(u)
}

func GetUserByID(tx *pop.Connection, id uuid.UUID) (*User, error) {
	u := User{}
	err := tx.Find(&u, id)

	if err != nil {
		return nil, fmt.Errorf("could not find user %v", err)
	}

	return &u, nil
}

func GetUserByUsername(tx *pop.Connection, username string) (*User, error) {
	u := User{}
	query := tx.Where("username = ?", username)
	err := query.First(&u)

	if err != nil {
		return nil, fmt.Errorf("could not find user %v", err)
	}

	return &u, nil
}

// Users is not required by pop and may be deleted
type Users []User

// String is not required by pop and may be deleted
func (u Users) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	var err error
	return validate.Validate(
		&validators.FuncValidator{
			Field:   u.Email,
			Name:    "Email",
			Message: "a user with email %s already exists",
			Fn: func() bool {
				var b bool
				q := tx.Where("email = ?", u.Email)
				b, err = q.Exists(u)
				if err != nil {
					return false
				}
				return !b
			},
		},
		&validators.FuncValidator{
			Field:   u.Email,
			Name:    "Email",
			Message: "email is empty",
			Fn: func() bool {
				return u.Email != ""
			},
		},
		&validators.FuncValidator{
			Field:   u.Username,
			Name:    "Username",
			Message: "username is empty",
			Fn: func() bool {
				return u.Username != ""
			},
		},
		&validators.FuncValidator{
			Field:   u.Username,
			Name:    "Username",
			Message: "username %s is already in use",
			Fn: func() bool {
				var b bool
				q := tx.Where("username = ?", u.Username)
				b, err = q.Exists(u)
				if err != nil {
					return false
				}
				return !b
			},
		},
		&validators.FuncValidator{
			Field:   u.UserType,
			Name:    "Role",
			Message: "no user type specified",
			Fn: func() bool {
				return u.UserType != ""
			},
		},
		&validators.FuncValidator{
			Field:   u.Username,
			Name:    "Username",
			Message: "user type must be 'artist' or 'follower'",
			Fn: func() bool {
				return u.UserType == "artist" || u.UserType == "follower"
			},
		},
	), err
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (u *User) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (u *User) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
