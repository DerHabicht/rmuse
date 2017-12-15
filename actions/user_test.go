package actions

/***********************************************************************************************************************
 * Login tests
 **********************************************************************************************************************/
// Test_Bad_User_Login attempts to login as an invalid user.
func (as *ActionSuite) Test_Bad_User_Login() {
	u := make(map[string]string)

	u["email"] = "dog@example.com"
	u["password"] = "goodpassword"

	res := as.JSON("/api/login").Post(u)
	as.Equal(401, res.Code)
}

// Test_Bad_Password_Login attempts to login as a valid user with an invalid password.
func (as *ActionSuite) Test_Bad_Password_Login() {
	u := make(map[string]string)

	u["email"] = "cat@example.com"
	u["password"] = "badpassword"

	res := as.JSON("/api/login").Post(u)
	as.Equal(401, res.Code)
}

// Test_Good_User_Login attempts to login as a valid user with a good password.
func (as *ActionSuite) Test_Good_User_Login() {
	u := make(map[string]string)

	u["email"] = "cat@example.com"
	u["password"] = "goodpassword"

	res := as.JSON("/api/login").Post(u)
	as.Contains("session_id", res.Code)
	as.Equal(200, res.Code)
}

// Test_User_Create attempts to create a new user.
func (as *ActionSuite) Test_User_Create() {
	as.Fail("Not Implemented!")
}

// Test_Duplicate_User_Create attempts to create a new user already in the database.
func (as *ActionSuite) Test_Duplicate_User_Create() {
	as.Fail("Not Implemented!")
}

// Test_Duplicate_Email_Create attempts to create a user with a different username but reuses an email address.
func (as *ActionSuite) Test_Duplicate_Email_Create() {
	as.Fail("Not Implemented!")
}
