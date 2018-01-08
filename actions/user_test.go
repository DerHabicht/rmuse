package actions

import (
	"net/http"

	"github.com/derhabicht/rmuse/models"
	"golang.org/x/crypto/bcrypt"
)

func (as *ActionSuite) Test_User_Create() {
	arg := struct{
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		Email     string `json:"email"`
		Username  string `json:"username"`
		Password  string `json:"password"`
	}{
		FirstName: "Oreo",
		LastName:  "Hawk",
		Email:     "cat@example.com",
		Username:  "oreo",
		Password:  "goodpassword",
	}

	res := as.JSON("/api/1/user").Post(arg)
	as.Equal(http.StatusOK, res.Code)
	as.Contains(res.Body.String(), "token")
	as.Contains(res.Body.String(), "Oreo")
	as.Contains(res.Body.String(), "Hawk")
	as.Contains(res.Body.String(), "cat@example.com")
	as.Contains(res.Body.String(), "oreo")
	as.NotContains(res.Body.String(), "password")

	u, err := models.GetUserByUsername(as.DB, "oreo")
	as.NoError(err)

	as.DB.Destroy(u)
	as.NoError(err)
}

func (as *ActionSuite) Test_User_Empty_Username_Create() {
	arg := struct{
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		Email     string `json:"email"`
		Username  string `json:"username"`
		Password  string `json:"password"`
	}{
		FirstName: "Oreo",
		LastName:  "Hawk",
		Email:     "blackcat@example.com",
		Username:  "",
		Password:  "goodpassword",
	}

	res := as.JSON("/api/1/user").Post(arg)
	as.Equal(http.StatusUnprocessableEntity, res.Code)
	as.Contains(res.Body.String(),"username is empty")
}

func (as *ActionSuite) Test_User_Empty_Email_Create() {
	arg := struct{
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		Email     string `json:"email"`
		Username  string `json:"username"`
		Password  string `json:"password"`
	}{
		FirstName: "Oreo",
		LastName:  "Hawk",
		Email:     "",
		Username:  "oreo",
		Password:  "goodpassword",
	}

	res := as.JSON("/api/1/user").Post(arg)
	as.Equal(http.StatusUnprocessableEntity, res.Code)
	as.Contains(res.Body.String(),"email is empty")
}

func (as *ActionSuite) Test_User_Duplicate_Username_Create() {
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
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		Email     string `json:"email"`
		Username  string `json:"username"`
		Password  string `json:"password"`
	}{
		FirstName: "Oreo",
		LastName:  "Hawk",
		Email:     "blackcat@example.com",
		Username:  "oreo",
		Password:  "goodpassword",
	}

	res := as.JSON("/api/1/user").Post(arg)
	as.Equal(http.StatusUnprocessableEntity, res.Code)
	as.Contains(res.Body.String(),"username oreo is already in use")

	err = as.DB.Destroy(&u)
	as.NoError(err)
}

func (as *ActionSuite) Test_User_Duplicate_Email_Create() {
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
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		Email     string `json:"email"`
		Username  string `json:"username"`
		Password  string `json:"password"`
	}{
		FirstName: "Oreo",
		LastName:  "Hawk",
		Email:     "cat@example.com",
		Username:  "cat",
		Password:  "goodpassword",
	}

	res := as.JSON("/api/1/user").Post(arg)
	as.Equal(http.StatusUnprocessableEntity, res.Code)
	as.Contains(res.Body.String(),"a user with email cat@example.com already exists")

	err = as.DB.Destroy(&u)
	as.NoError(err)
}
