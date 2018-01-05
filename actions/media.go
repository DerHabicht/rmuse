package actions

import (
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/markbates/pop"
	"github.com/pkg/errors"

	"github.com/derhabicht/rmuse/models"
)

// MediaGet default implementation.
func MediaGet(c buffalo.Context) error {
	return c.Render(200, r.HTML("media/get.html"))
}

// MediaUpload default implementation.
func MediaUpload(c buffalo.Context) error {
	u, ok := c.Value("user").(*models.User)

	if !ok || u == nil {
		return c.Error(http.StatusUnauthorized, fmt.Errorf("must be logged in to upload media"))
	}

	m := &models.Medium{}

	if err := c.Bind(m); err != nil {
		return errors.WithStack(err)
	}

	m.User = u.ID

	if m.Permission == "" {
		m.Permission = "public"
	}

	tx := c.Value("tx").(*pop.Connection)
	verrs, err := m.Create(tx)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
	}

	return c.Render(http.StatusOK, nil)
}
