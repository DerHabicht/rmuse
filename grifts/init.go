package grifts

import (
	"github.com/derhabicht/rmuse/actions"
	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
