package vari

import (
	"github.com/easysoft/zendata/src/model"
	"github.com/easysoft/zendata/src/utils/const"
)

var (
	Config      = model.Config{Version: 1, Language: "en"}

	RunMode      constant.RunMode

	ExeDir       string

	LogDir       string
	ScreenWidth  int
	ScreenHeight int

	RequestType  string
	Verbose     bool
	Interpreter string

	WithHead bool
	HeadSep string

	JsonResp string = "[]"
	Ip string
	Port int

	Def = model.DefData{}
	Res = map[string]map[string][]string{}

	DefaultDir string
	ConfigDir string
)
