package actions

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gobuffalo/buffalo"
	"github.com/markbates/pop"
	"github.com/pkg/errors"

	"github.com/derhabicht/rmuse/models"
	"github.com/satori/go.uuid"
)

// MediaGet default implementation.
func MediaGet(c buffalo.Context) error {
	var media []*models.Medium
	tx := c.Value("tx").(*pop.Connection)

	if p, ok := c.Params().(url.Values)["id"]; ok {
		for _, uuidStr := range p {
			uuid, err := uuid.FromString(uuidStr)
			if err == nil {
				m, err := models.GetMediumByID(tx, uuid)
				if err == nil {
					media = append(media, m)
				}
			}
		}
		return c.Render(http.StatusOK, r.JSON(media))
	}

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
