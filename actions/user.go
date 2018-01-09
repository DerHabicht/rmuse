package actions

import (
	"net/http"

	"github.com/derhabicht/rmuse/models"
	"github.com/gobuffalo/buffalo"
	"github.com/markbates/pop"
	"golang.org/x/crypto/bcrypt"
)

// UserCreate default implementation.
func UserCreate(c buffalo.Context) error {
	type argument struct{
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		Email     string `json:"email"`
		Username  string `json:"username"`
		UserType  string `json:"type"`
		Password  string `json:"password"`
	}

	arg := &argument{}
	if err := c.Bind(arg); err != nil {
		return c.Render(http.StatusUnprocessableEntity, r.JSON("{\"error\":\"malformed argument body\"}"))
	}

	ph, err := bcrypt.GenerateFromPassword([]byte(arg.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Render(http.StatusInternalServerError, r.JSON("{\"error\":\"cannot hash password\"}"))
	}

	u := &models.User {
		FirstName:    arg.FirstName,
		LastName:     arg.LastName,
		Email:        arg.Email,
		Username:     arg.Username,
		UserType:     arg.UserType,
		PasswordHash: string(ph),
	}

	tx := c.Value("tx").(*pop.Connection)
	verrs, err := u.Create(tx)
	if err != nil {
		// TODO: Double check validations here to see why they fail
		return c.Render(http.StatusInternalServerError, r.JSON("{\"error\":\"failed to create user\"}"))
	}

	if verrs.HasAny() {
		return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
	}

	ts, err := u.CreateJWTToken()
	if err != nil {
		return c.Render(http.StatusInternalServerError, r.JSON("{\"error\":\"failed to create token\"}"))
	}

	res := struct {
		Token string       `json:"token"`
		User *models.User  `json:"user"`
	}{
		Token: ts,
		User:  u,
	}

	return c.Render(http.StatusOK, r.JSON(res))
}

func UserRead(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.JSON(c.Value("user")))
}
