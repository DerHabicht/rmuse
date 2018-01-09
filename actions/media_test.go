package actions

import (
	"fmt"
	"net/http"

	"github.com/derhabicht/rmuse/models"
	"golang.org/x/crypto/bcrypt"
)

func (as *ActionSuite) Test_Media_Upload_Logged_Out() {
	req := as.JSON("/api/1/media")
	req.Headers["Authorization"] = ""

	arg := struct {
		URI      string `json:"uri"`
		FileType string `json:"type"`
	}{
		URI:      "someplace",
		FileType: "image/png",
	}
	res := req.Post(arg)

	as.Equal(http.StatusUnauthorized, res.Code)

	as.DB.RawQuery("DELETE FROM users")
	as.DB.RawQuery("DELETE FROM media")
}

func (as *ActionSuite) Test_Media_Upload_Logged_In() {
	ph, err := bcrypt.GenerateFromPassword([]byte("goodpassword"), bcrypt.DefaultCost)
	as.NoError(err)

	user := models.User{
		FirstName:    "Oreo",
		LastName:     "Hawk",
		Email:        "cat@example.com",
		Username:     "oreo",
		UserType:     "artist",
		PasswordHash: string(ph),
	}

	err = as.DB.Create(&user)
	as.NoError(err)

	u, err := models.GetUserByUsername(as.DB, "oreo")

	req := as.JSON("/api/1/media")
	req.Headers["Authorization"], err = u.CreateJWTToken()
	as.NoError(err)

	arg := struct {
		URI      string `json:"uri"`
		FileType string `json:"type"`
	}{
		URI:      "someplace",
		FileType: "image/png",
	}
	res := req.Post(arg)

	as.Equal(http.StatusOK, res.Code)

	as.DB.RawQuery("DELETE FROM users")
	as.DB.RawQuery("DELETE FROM media")
}

func (as *ActionSuite) Test_Media_Upload_Empty_Type() {
	ph, err := bcrypt.GenerateFromPassword([]byte("goodpassword"), bcrypt.DefaultCost)
	as.NoError(err)

	user := models.User{
		FirstName:    "Oreo",
		LastName:     "Hawk",
		Email:        "cat@example.com",
		Username:     "oreo",
		PasswordHash: string(ph),
		UserType:     "artist",
	}

	err = as.DB.Create(&user)
	as.NoError(err)

	u, err := models.GetUserByUsername(as.DB, "oreo")

	req := as.JSON("/api/1/media")
	req.Headers["Authorization"], err = u.CreateJWTToken()
	as.NoError(err)

	arg := struct {
		URI      string `json:"uri"`
		FileType string `json:"type"`
	}{
		URI:      "someplace",
		FileType: "",
	}
	res := req.Post(arg)

	as.Equal(http.StatusUnprocessableEntity, res.Code)
	as.Contains(res.Body.String(), "type is empty")

	as.DB.RawQuery("DELETE FROM users")
	as.DB.RawQuery("DELETE FROM media")
}

func (as *ActionSuite) Test_Media_Upload_Empty_URI() {
	ph, err := bcrypt.GenerateFromPassword([]byte("goodpassword"), bcrypt.DefaultCost)
	as.NoError(err)

	user := models.User{
		FirstName:    "Oreo",
		LastName:     "Hawk",
		Email:        "cat@example.com",
		Username:     "oreo",
		PasswordHash: string(ph),
		UserType:     "artist",
	}

	err = as.DB.Create(&user)
	as.NoError(err)

	u, err := models.GetUserByUsername(as.DB, "oreo")

	req := as.JSON("/api/1/media")
	req.Headers["Authorization"], err = u.CreateJWTToken()
	as.NoError(err)

	arg := struct {
		URI      string `json:"uri"`
		FileType string `json:"type"`
	}{
		URI:      "",
		FileType: "image/png",
	}
	res := req.Post(arg)

	as.Equal(http.StatusUnprocessableEntity, res.Code)
	as.Contains(res.Body.String(), "uri is empty")

	as.DB.RawQuery("DELETE FROM users")
	as.DB.RawQuery("DELETE FROM media")
}

func (as *ActionSuite) Test_Media_Upload_Follower() {
	ph, err := bcrypt.GenerateFromPassword([]byte("goodpassword"), bcrypt.DefaultCost)
	as.NoError(err)

	user := models.User{
		FirstName:    "Oreo",
		LastName:     "Hawk",
		Email:        "cat@example.com",
		Username:     "oreo",
		UserType:     "follower",
		PasswordHash: string(ph),
	}

	err = as.DB.Create(&user)
	as.NoError(err)

	u, err := models.GetUserByUsername(as.DB, "oreo")

	req := as.JSON("/api/1/media")
	req.Headers["Authorization"], err = u.CreateJWTToken()
	as.NoError(err)

	arg := struct {
		URI      string `json:"uri"`
		FileType string `json:"type"`
	}{
		URI:      "someplace",
		FileType: "image/png",
	}
	res := req.Post(arg)

	as.Equal(http.StatusUnauthorized, res.Code)
	as.Contains(res.Body.String(), "must be artist to upload media")

	as.DB.RawQuery("DELETE FROM users")
	as.DB.RawQuery("DELETE FROM media")
}

func (as *ActionSuite) Test_Media_Upload_Duplicate_URI() {
	ph, err := bcrypt.GenerateFromPassword([]byte("goodpassword"), bcrypt.DefaultCost)
	as.NoError(err)

	user := models.User{
		FirstName:    "Oreo",
		LastName:     "Hawk",
		Email:        "cat@example.com",
		Username:     "oreo",
		PasswordHash: string(ph),
		UserType:     "artist",
	}

	err = as.DB.Create(&user)
	as.NoError(err)

	medium := models.Medium{
		URI:      "someplace",
		Filetype: "image/png",
	}

	err = as.DB.Create(&medium)
	as.NoError(err)

	u, err := models.GetUserByUsername(as.DB, "oreo")

	req := as.JSON("/api/1/media")
	req.Headers["Authorization"], err = u.CreateJWTToken()
	as.NoError(err)

	arg := struct {
		URI      string `json:"uri"`
		FileType string `json:"type"`
	}{
		URI:      "someplace",
		FileType: "image/png",
	}
	res := req.Post(arg)

	as.Equal(http.StatusUnprocessableEntity, res.Code)
	as.Contains(res.Body.String(), "there is already a file with URI someplace")

	as.DB.RawQuery("DELETE FROM users")
	as.DB.RawQuery("DELETE FROM media")
}

func (as *ActionSuite) Test_Media_Get_Public() {
	medium := models.Medium{
		URI:        "someplace",
		Filetype:   "image/png",
		Permission: "public",
	}

	err := as.DB.Create(&medium)
	as.NoError(err)

	id, err := models.GetMediumIDByURI(as.DB, "someplace")
	as.NoError(err)

	res := as.JSON(fmt.Sprintf("/api/1/media?id=%s", id.String())).Get()

	as.Equal(http.StatusOK, res.Code)

	as.DB.RawQuery("DELETE FROM media")
}

func (as *ActionSuite) Test_Media_Get_Follower() {
	ph, err := bcrypt.GenerateFromPassword([]byte("goodpassword"), bcrypt.DefaultCost)
	as.NoError(err)

	user := models.User{
		FirstName:    "Oreo",
		LastName:     "Hawk",
		Email:        "cat@example.com",
		Username:     "oreo",
		PasswordHash: string(ph),
		UserType:     "follower",
	}

	err = as.DB.Create(&user)
	as.NoError(err)

	ph, err = bcrypt.GenerateFromPassword([]byte("goodpassword"), bcrypt.DefaultCost)
	as.NoError(err)

	user = models.User{
		FirstName:    "Raja",
		LastName:     "Hawk",
		Email:        "clutz@example.com",
		Username:     "raja",
		PasswordHash: string(ph),
		UserType:     "artist",
	}

	err = as.DB.Create(&user)
	as.NoError(err)

	raj, err := models.GetUserByUsername(as.DB, "raja")
	oreo, err := models.GetUserByUsername(as.DB, "oreo")

	f := models.Follow{
		Follower: oreo.ID,
		Followed: raj.ID,
	}

	err = as.DB.Create(&f)
	as.NoError(err)

	medium := models.Medium{
		URI:        "someplace",
		Filetype:   "image/png",
		User:       raj.ID,
		Permission: "follower",
	}

	err = as.DB.Create(&medium)
	as.NoError(err)

	id, err := models.GetMediumIDByURI(as.DB, "someplace")
	as.NoError(err)

	req := as.JSON(fmt.Sprintf("/api/1/media?id=%s", id.String()))
	req.Headers["Authorization"], err = oreo.CreateJWTToken()
	as.NoError(err)

	res := req.Get()

	as.Equal(http.StatusOK, res.Code)

	as.DB.RawQuery("DELETE FROM users")
	as.DB.RawQuery("DELETE FROM media")
	as.DB.RawQuery("DELETE FROM follows")
}

func (as *ActionSuite) Test_Media_Get_Public_Unauthorized() {
	medium := models.Medium{
		URI:        "someplace",
		Filetype:   "imapge/png",
		Permission: "follower",
	}

	err := as.DB.Create(&medium)
	as.NoError(err)

	id, err := models.GetMediumIDByURI(as.DB, "someplace")
	as.NoError(err)

	res := as.JSON(fmt.Sprintf("/api/1/media?id=%s", id.String())).Get()

	as.Equal(http.StatusUnauthorized, res.Code)

	as.DB.RawQuery("DELETE FROM media")
}

func (as *ActionSuite) Test_Media_Get_Follower_Unauthorized() {
	ph, err := bcrypt.GenerateFromPassword([]byte("goodpassword"), bcrypt.DefaultCost)
	as.NoError(err)

	user := models.User{
		FirstName:    "Oreo",
		LastName:     "Hawk",
		Email:        "cat@example.com",
		Username:     "oreo",
		PasswordHash: string(ph),
		UserType:     "follower",
	}

	err = as.DB.Create(&user)
	as.NoError(err)

	ph, err = bcrypt.GenerateFromPassword([]byte("goodpassword"), bcrypt.DefaultCost)
	as.NoError(err)

	user = models.User{
		FirstName:    "Raja",
		LastName:     "Hawk",
		Email:        "clutz@example.com",
		Username:     "raja",
		PasswordHash: string(ph),
		UserType:     "artist",
	}

	err = as.DB.Create(&user)
	as.NoError(err)

	raj, err := models.GetUserByUsername(as.DB, "raja")
	oreo, err := models.GetUserByUsername(as.DB, "oreo")

	medium := models.Medium{
		URI:        "someplace",
		Filetype:   "image/png",
		User:       raj.ID,
		Permission: "follower",
	}

	err = as.DB.Create(&medium)
	as.NoError(err)

	id, err := models.GetMediumIDByURI(as.DB, "someplace")
	as.NoError(err)

	req := as.JSON(fmt.Sprintf("/api/1/media?id=%s", id.String()))
	req.Headers["Authorization"], err = oreo.CreateJWTToken()
	as.NoError(err)

	res := req.Get()

	as.Equal(http.StatusUnauthorized, res.Code)

	as.DB.RawQuery("DELETE FROM users")
	as.DB.RawQuery("DELETE FROM media")
}
