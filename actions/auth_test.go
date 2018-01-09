package actions

import (
	"net/http"
	"golang.org/x/crypto/bcrypt"

	"github.com/derhabicht/rmuse/models"
)

// Test_Bad_User_Login attempts to login as an invalid user.
func (as *ActionSuite) Test_Login_Bad_User() {
	arg := struct{
		Email     string `json:"email"`
		Password  string `json:"password"`
	}{
		Email:    "dog@example.com",
		Password: "goodpassword",
	}

	res := as.JSON("/api/1/login").Post(arg)
	as.Equal(http.StatusUnprocessableEntity, res.Code)
	as.Contains(res.Body.String(), "invalid email or password")
}

// Test_Bad_Password_Login attempts to login as a valid user with an invalid password.
func (as *ActionSuite) Test_Login_Bad_Password() {
	ph, err := bcrypt.GenerateFromPassword([]byte("goodpassword"), bcrypt.DefaultCost)
	as.NoError(err)

	u := models.User {
		FirstName:    "Oreo",
		LastName:     "Hawk",
		Email:        "cat@example.com",
		Username:     "oreo",
		UserType:     "artist",
		PasswordHash: string(ph),
	}

	err = as.DB.Create(&u)
	as.NoError(err)

	arg := struct{
		Email     string `json:"email"`
		Password  string `json:"password"`
	}{
		Email:    "cat@example.com",
		Password: "badpassword",
	}

	res := as.JSON("/api/1/login").Post(arg)
	as.Equal(http.StatusUnprocessableEntity, res.Code)
	as.Contains(res.Body.String(), "invalid email or password")

	as.DB.RawQuery("DELETE FROM users")
}

// Test_Good_User_Login attempts to login as a valid user with a good password.
func (as *ActionSuite) Test_Login_Good_User() {
	ph, err := bcrypt.GenerateFromPassword([]byte("goodpassword"), bcrypt.DefaultCost)
	as.NoError(err)

	u := models.User {
		FirstName:    "Oreo",
		LastName:     "Hawk",
		Email:        "cat@example.com",
		Username:     "oreo",
		PasswordHash: string(ph),
		UserType:     "artist",
	}

	err = as.DB.Create(&u)
	as.NoError(err)

	arg := struct{
		Email     string `json:"email"`
		Password  string `json:"password"`
	}{
		Email:    "cat@example.com",
		Password: "goodpassword",
	}

	res := as.JSON("/api/1/login").Post(arg)
	as.Equal(http.StatusOK, res.Code)
	as.Contains(res.Body.String(), "token")
	as.Contains(res.Body.String(), "Oreo")
	as.Contains(res.Body.String(), "Hawk")
	as.Contains(res.Body.String(), "cat@example.com")
	as.Contains(res.Body.String(), "oreo")
	as.NotContains(res.Body.String(), "password")

	as.DB.RawQuery("DELETE FROM users")
}


