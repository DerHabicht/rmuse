package actions

import (
	"testing"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/suite"
)

type ActionSuite struct {
	*suite.Action
}

func Test_ActionSuite(t *testing.T) {
	GOPATH := envy.Get("GOPATH", "")
	envy.Set("JWT_KEY_PATH", GOPATH+"/src/github.com/derhabicht/rmuse/jwtRS256.key")
	as := &ActionSuite{suite.NewAction(App())}
	suite.Run(t, as)
}
