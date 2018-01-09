package actions

import (
	"github.com/pkg/errors"
	"github.com/derhabicht/rmuse/models"
	"github.com/gobuffalo/buffalo"
	"github.com/markbates/pop"
	"net/http"
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
		return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
	}

	// Now that the user has been created, we need to drop the password from the user struct so we don't end up
	// sending the password back in our response later.
	u.Password = ""

	ts, err := u.CreateJWTToken()
	if err != nil {
		return errors.WithStack(err)
	}

	res := struct {
		Token string `json:"token"`
	}{
		ts,
	}

	return c.Render(http.StatusOK, r.JSON(res))
}

func UserRead(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.JSON(c.Value("user")))
}
