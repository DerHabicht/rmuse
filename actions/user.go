package actions

import (
	"net/http"

	"github.com/derhabicht/rmuse/models"
	"github.com/gobuffalo/buffalo"
	"github.com/markbates/pop"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// UserCreate default implementation.
func UserCreate(c buffalo.Context) error {
	type argument struct{
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		Email     string `json:"email"`
		Username  string `json:"username"`
		Password  string `json:"password"`
	}

	arg := &argument{}
	if err := c.Bind(arg); err != nil {
		return errors.WithStack(err)
	}

	ph, err := bcrypt.GenerateFromPassword([]byte(arg.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.WithStack(err)
	}

	u := &models.User {
		FirstName:    arg.FirstName,
		LastName:     arg.LastName,
		Email:        arg.Email,
		Username:     arg.Username,
		PasswordHash: string(ph),
	}

	tx := c.Value("tx").(*pop.Connection)
	verrs, err := u.Create(tx)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
	}

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
