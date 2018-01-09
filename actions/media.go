package actions

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gobuffalo/buffalo"
	"github.com/markbates/pop"
	"github.com/satori/go.uuid"

	"github.com/derhabicht/rmuse/models"
)

// MediaGet default implementation.
func MediaGet(c buffalo.Context) error {
	u, ok := c.Value("user").(*models.User)

	if !ok {
		u = nil
	}

	var media []*models.Medium
	tx := c.Value("tx").(*pop.Connection)

	if p, ok := c.Params().(url.Values)["id"]; ok {
		for _, uuidStr := range p {
			uuid, err := uuid.FromString(uuidStr)
			if err == nil {
				m, err := models.GetMediumByID(tx, uuid, u)
				if err == nil {
					media = append(media, m)
				}
			}
		}
		return c.Render(http.StatusOK, r.JSON(media))
	}

	return c.Render(http.StatusUnauthorized, nil)
}

// MediaUpload default implementation.
func MediaUpload(c buffalo.Context) error {
	u, ok := c.Value("user").(*models.User)

	if !ok || u == nil {
		return c.Render(http.StatusUnauthorized, r.JSON("must be logged in to upload files"))
	}


	m := &models.Medium{}

	if err := c.Bind(m); err != nil {
		return c.Error(http.StatusInternalServerError, fmt.Errorf("unable to bind medium %v", err))
	}

	m.User = u.ID

	if m.Permission == "" {
		m.Permission = "public"
	}

	tx := c.Value("tx").(*pop.Connection)
	verrs, err := m.Create(tx)
	if err != nil {
		return c.Error(http.StatusInternalServerError, fmt.Errorf("unable to create medium %v", err))
	}

	if verrs.HasAny() {
		return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
	}

	return c.Render(http.StatusOK, r.JSON(m))
}
