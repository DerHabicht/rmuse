package actions

import (
	"github.com/pkg/errors"
	"github.com/derhabicht/rmuse/models"
	"github.com/gobuffalo/buffalo"
	"github.com/markbates/pop"
)

// UserCreate default implementation.
func UserCreate(c buffalo.Context) error {
	u := &models.User{}

	if err := c.Bind(u); err != nil {
		return errors.WithStack(err)
	}

	tx := c.Value("tx").(*pop.Connection)
	verrs, err := u.Create(tx)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		return c.Render(422, r.JSON(verrs))
	}

	// Now that the user has been created, we need to drop the password from the user struct so we don't end up
	// sending the password back in our response later.
	u.Password = ""

	ts, err := u.CreateJWTToken()
	if err != nil {
		return errors.WithStack(err)
	}

	res := struct {
		token string
	}{
		ts,
	}

	return c.Render(200, r.JSON(res))
}
