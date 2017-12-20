package actions

import (
	"strings"
	"database/sql"

	"github.com/gobuffalo/buffalo"
	"github.com/markbates/pop"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/derhabicht/rmuse/models"
)


func AuthCreateSession(c buffalo.Context) error {
	bad := func() error {
		m := struct {
			Message string `json:"error"`
		}{
			"invalid email/password",
		}
		return c.Render(422, r.JSON(m))
	}

	u := &models.User{}
	if err := c.Bind(u); err != nil {
		return errors.WithStack(err)
	}

	tx := c.Value("tx").(*pop.Connection)

	// Try to find the user by email
	err := tx.Where("email = ?", strings.ToLower(u.Email)).First(u)

	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return bad()
		}
		return errors.WithStack(err)
	}

	// Test the user's password
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(u.Password))
	if err != nil {
		return bad()
	}

	// Now that the user has been authenticated, we need to drop the password from the user struct so we don't end up
	// sending the password back in our response later.
	u.Password = ""

	ts, err := u.CreateJWTToken()
	if err != nil {
		return errors.WithStack(err)
	}

	res := struct {
		Token string    `json:"token"`
		Username string `json:"username"`
	}{
		ts,
		u.Username,
	}

	return c.Render(200, r.JSON(res))
}

