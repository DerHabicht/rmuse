package actions

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/markbates/pop"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"github.com/derhabicht/rmuse/models"
	"github.com/satori/go.uuid"
)

func AuthCreateSession(c buffalo.Context) error {
	bad := func() error {
		return c.Error(http.StatusUnprocessableEntity, fmt.Errorf("invalid email or password"))
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

	return c.Render(http.StatusOK, r.JSON(res))
}

func VerifyToken(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")

		if len(tokenString) == 0 {
			return c.Error(http.StatusUnauthorized, fmt.Errorf("no token set in headers"))
		}

		// parsing token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// key
			sk, err := ioutil.ReadFile(envy.Get("JWT_KEY_PATH", "jwtRS256.key"))

			if err != nil {
				return nil, fmt.Errorf("could not open jwt key, %v", err)
			}

			return sk, nil
		})

		if err != nil {
			return c.Error(http.StatusUnauthorized, fmt.Errorf("could not parse the token, %v", err))
		}

		// getting claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			tx := c.Value("tx").(*pop.Connection)

			logrus.Errorf("claims: %v", claims)

			// retrieving user from db
			id, err := uuid.FromString(claims["jti"].(string))
			if err != nil {
				return c.Error(http.StatusUnauthorized, fmt.Errorf("could not identify the user"))
			}

			u, err := models.GetUserByID(tx, id)

			if err != nil {
				return c.Error(http.StatusUnauthorized, fmt.Errorf("could not identify the user"))
			}

			c.Set("user", u)

		} else {
			return c.Error(http.StatusUnauthorized, fmt.Errorf("failed to validate token: %v", claims))
		}

		return next(c)
	}
}
