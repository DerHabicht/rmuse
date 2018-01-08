package actions

import (
	"net/http"
	"golang.org/x/crypto/bcrypt"

	"github.com/derhabicht/rmuse/models"
)

// Test_Bad_User_Login attempts to login as an invalid user.
func (as *ActionSuite) Test_Bad_User_Login() {
	arg := struct{
		Email     string `json:"email"`
		Password  string `json:"password"`
	}{
		Email:    "dog@example.com",
		Password: "goodpassword",
	}

	res := as.JSON("/api/1/login").Post(arg)
	as.Equal(http.StatusUnprocessableEntity, res.Code)
}

// Test_Bad_Password_Login attempts to login as a valid user with an invalid password.
func (as *ActionSuite) Test_Bad_Password_Login() {
	arg := struct{
		Email     string `json:"email"`
		Password  string `json:"password"`
	}{
		Email:    "cat@example.com",
		Password: "badpassword",
	}

	res := as.JSON("/api/1/login").Post(arg)
	as.Equal(http.StatusUnprocessableEntity, res.Code)
}

// Test_Good_User_Login attempts to login as a valid user with a good password.
func (as *ActionSuite) Test_Good_User_Login() {
	ph, err := bcrypt.GenerateFromPassword([]byte("goodpassword"), bcrypt.DefaultCost)
	as.NoError(err)

	u := models.User {
		FirstName:    "Oreo",
		LastName:     "Hawk",
		Email:        "cat@example.com",
		Username:     "oreo",
		PasswordHash: string(ph),
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
	as.Contains(res.Body.String(), "token")
	as.Equal(http.StatusOK, res.Code)

	err = as.DB.Destroy(&u)
	as.NoError(err)
}


