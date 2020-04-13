package vari

import (
	"github.com/easysoft/zendata/src/model"
	"github.com/easysoft/zendata/src/utils/const"
)

var (
	Config         = model.Config{}

	RunMode     constant.RunMode
	ZDataDir    string
	LogDir      string
	ScreenWidth int
	ScreenHeight int

	RequestType  string
	Verbose     bool
	Interpreter string
)
