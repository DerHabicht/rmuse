package actions

import (
	"fmt"
	"net/http"

	"github.com/derhabicht/rmuse/models"
	"golang.org/x/crypto/bcrypt"
)

func (as *ActionSuite) Test_User_Create() {
	arg := struct {
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		Email     string `json:"email"`
		Username  string `json:"username"`
		Password  string `json:"password"`
		UserType  string `json:"type"`
	}{
		FirstName: "Oreo",
		LastName:  "Hawk",
		Email:     "cat@example.com",
		Username:  "oreo",
		Password:  "goodpassword",
		UserType:  "artist",
	}

	res := as.JSON("/api/1/user").Post(arg)
	as.Equal(http.StatusOK, res.Code)
	as.Contains(res.Body.String(), "token")
	as.Contains(res.Body.String(), "Oreo")
	as.Contains(res.Body.String(), "Hawk")
	as.Contains(res.Body.String(), "cat@example.com")
	as.Contains(res.Body.String(), "oreo")
	as.NotContains(res.Body.String(), "password")

	as.DB.RawQuery("DELETE FROM users")
}

func (as *ActionSuite) Test_User_Empty_Username_Create() {
	arg := struct {
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		Email     string `json:"email"`
		Username  string `json:"username"`
		Password  string `json:"password"`
		UserType  string `json:"type"`
	}{
		FirstName: "Oreo",
		LastName:  "Hawk",
		Email:     "blackcat@example.com",
		Username:  "",
		Password:  "goodpassword",
		UserType:  "artist",
	}

	res := as.JSON("/api/1/user").Post(arg)
	as.Equal(http.StatusUnprocessableEntity, res.Code)
	as.Contains(res.Body.String(), "username is empty")
}

func (as *ActionSuite) Test_User_Empty_Email_Create() {
	arg := struct {
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		Email     string `json:"email"`
		Username  string `json:"username"`
		Password  string `json:"password"`
		UserType  string `json:"type"`
	}{
		FirstName: "Oreo",
		LastName:  "Hawk",
		Email:     "",
		Username:  "oreo",
		Password:  "goodpassword",
		UserType:  "artist",
	}

	res := as.JSON("/api/1/user").Post(arg)
	as.Equal(http.StatusUnprocessableEntity, res.Code)
	as.Contains(res.Body.String(), "email is empty")
}

func (as *ActionSuite) Test_User_Empty_Type_Create() {
	arg := struct {
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		Email     string `json:"email"`
		Username  string `json:"username"`
		Password  string `json:"password"`
		UserType  string `json:"type"`
	}{
		FirstName: "Oreo",
		LastName:  "Hawk",
		Email:     "blackcat@example.com",
		Username:  "oreo",
		Password:  "goodpassword",
		UserType:  "",
	}

	res := as.JSON("/api/1/user").Post(arg)
	as.Equal(http.StatusUnprocessableEntity, res.Code)
	as.Contains(res.Body.String(), "no user type specified")
}

func (as *ActionSuite) Test_User_Bad_Type_Create() {
	arg := struct {
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		Email     string `json:"email"`
		Username  string `json:"username"`
		Password  string `json:"password"`
		UserType  string `json:"type"`
	}{
		FirstName: "Oreo",
		LastName:  "Hawk",
		Email:     "blackcat@example.com",
		Username:  "oreo",
		Password:  "goodpassword",
		UserType:  "cat",
	}

	res := as.JSON("/api/1/user").Post(arg)
	as.Equal(http.StatusUnprocessableEntity, res.Code)
	as.Contains(res.Body.String(), "user type must be 'artist' or 'follower'")
}

func (as *ActionSuite) Test_User_Duplicate_Username_Create() {
	ph, err := bcrypt.GenerateFromPassword([]byte("goodpassword"), bcrypt.DefaultCost)
	as.NoError(err)

	u := models.User{
		FirstName:    "Oreo",
		LastName:     "Hawk",
		Email:        "cat@example.com",
		Username:     "oreo",
		PasswordHash: string(ph),
		UserType:     "artist",
	}

	err = as.DB.Create(&u)
	as.NoError(err)

	arg := struct {
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		Email     string `json:"email"`
		Username  string `json:"username"`
		Password  string `json:"password"`
		UserType  string `json:"type"`
	}{
		FirstName: "Oreo",
		LastName:  "Hawk",
		Email:     "blackcat@example.com",
		Username:  "oreo",
		Password:  "goodpassword",
		UserType:  "artist",
	}

	res := as.JSON("/api/1/user").Post(arg)
	as.Equal(http.StatusUnprocessableEntity, res.Code)
	as.Contains(res.Body.String(), "username oreo is already in use")

	as.DB.RawQuery("DELETE FROM users")
}

func (as *ActionSuite) Test_User_Duplicate_Email_Create() {
	ph, err := bcrypt.GenerateFromPassword([]byte("goodpassword"), bcrypt.DefaultCost)
	as.NoError(err)

	u := models.User{
		FirstName:    "Oreo",
		LastName:     "Hawk",
		Email:        "cat@example.com",
		Username:     "oreo",
		PasswordHash: string(ph),
		UserType:     "artist",
	}

	err = as.DB.Create(&u)
	as.NoError(err)

	arg := struct {
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		Email     string `json:"email"`
		Username  string `json:"username"`
		Password  string `json:"password"`
		UserType  string `json:"type"`
	}{
		FirstName: "Oreo",
		LastName:  "Hawk",
		Email:     "cat@example.com",
		Username:  "cat",
		Password:  "goodpassword",
		UserType:  "artist",
	}

	res := as.JSON("/api/1/user").Post(arg)
	as.Equal(http.StatusUnprocessableEntity, res.Code)
	as.Contains(res.Body.String(), "a user with email cat@example.com already exists")

	as.DB.RawQuery("DELETE FROM users")
}

func (as *ActionSuite) Test_User_Page_Fetch() {
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

	m1 := models.Medium{
		User:     u.ID,
		URI:      "something1",
		Filetype: "image/png",
		PosX:     1,
		PosY:     1,
	}

	m2 := models.Medium{
		User:     u.ID,
		URI:      "something2",
		Filetype: "image/png",
		PosX:     2,
		PosY:     2,
	}

	m3 := models.Medium{
		User:     u.ID,
		URI:      "something3",
		Filetype: "image/png",
		PosX:     3,
		PosY:     3,
	}

	err = as.DB.Create(&m1)
	as.NoError(err)
	err = as.DB.Create(&m2)
	as.NoError(err)
	err = as.DB.Create(&m3)
	as.NoError(err)

	res := as.JSON("/api/1/user/oreo").Get()
	fmt.Println(res.Body.String())
	as.Equal(http.StatusOK, res.Code)
	as.Contains(res.Body.String(), "something1")
	as.Contains(res.Body.String(), "something2")
	as.Contains(res.Body.String(), "something3")

	as.DB.RawQuery("DELETE FROM users")
	as.DB.RawQuery("DELETE FROM media")
}

func (as *ActionSuite) Test_User_Update() {
	ph, err := bcrypt.GenerateFromPassword([]byte("goodpassword"), bcrypt.DefaultCost)
	as.NoError(err)

	u := models.User{
		FirstName:    "Oreo",
		LastName:     "Hawk",
		Email:        "cat@example.com",
		Username:     "oreo",
		PasswordHash: string(ph),
		UserType:     "artist",
	}

	err = as.DB.Create(&u)
	as.NoError(err)
	arg := struct {
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		Email     string `json:"email"`
		Username  string `json:"username"`
		UserType  string `json:"type"`
		Password  string `json:"password"`
	}{
		FirstName: "Raja",
		LastName:  "Hawk",
		Email:     "cat@example.com",
		Username:  "oreo",
		Password:  "goodpassword",
		UserType:  "artist",
	}

	req := as.JSON("/api/1/user")
	req.Headers["Authorization"], err = u.CreateJWTToken()
	as.NoError(err)
	res := req.Put(arg)

	as.Equal(http.StatusOK, res.Code)
	as.Contains(res.Body.String(), "Raja")
	as.Contains(res.Body.String(), "Hawk")
	as.Contains(res.Body.String(), "cat@example.com")
	as.Contains(res.Body.String(), "oreo")
	as.NotContains(res.Body.String(), "password")

	as.DB.RawQuery("DELETE FROM users")
}

func (as *ActionSuite) Test_User_Update_Bad_Password() {
	ph, err := bcrypt.GenerateFromPassword([]byte("goodpassword"), bcrypt.DefaultCost)
	as.NoError(err)

	u := models.User{
		FirstName:    "Oreo",
		LastName:     "Hawk",
		Email:        "cat@example.com",
		Username:     "oreo",
		PasswordHash: string(ph),
		UserType:     "artist",
	}

	err = as.DB.Create(&u)
	as.NoError(err)
	arg := struct {
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		Email     string `json:"email"`
		Username  string `json:"username"`
		Password  string `json:"password"`
		UserType  string `json:"type"`
	}{
		FirstName: "Raja",
		LastName:  "Hawk",
		Email:     "cat@example.com",
		Username:  "oreo",
		Password:  "badpassword",
		UserType:  "artist",
	}

	res := as.JSON("/api/1/user").Put(arg)
	as.Equal(http.StatusUnauthorized, res.Code)

	as.DB.RawQuery("DELETE FROM users")
}
