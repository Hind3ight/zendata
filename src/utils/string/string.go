package stringUtils

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	constant "github.com/easysoft/zendata/src/utils/const"
	"github.com/easysoft/zendata/src/utils/vari"
	"github.com/mattn/go-runewidth"
	"strconv"
	"strings"
	"unicode"
)

func TrimAll(str string) string {
	str = strings.Trim(str, "\n")
	str = strings.TrimSpace(str)

	return str
}

func Ucfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}
func UcAll(str string) string {
	ret := ""

	for _, v := range str {
		ret += string(unicode.ToUpper(v))
	}
	return ret
}

func Lcfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

func U2s(form string) (to string, err error) {
	bs, err := hex.DecodeString(strings.Replace(form, `\u`, ``, -1))
	if err != nil {
		return
	}
	for i, bl, br, r := 0, len(bs), bytes.NewReader(bs), uint16(0); i < bl; i += 2 {
		binary.Read(br, binary.BigEndian, &r)
		to += string(r)
	}
	return
}

func BoolToPass(b bool) string {
	if b {
		return constant.PASS.String()
	} else {
		return constant.FAIL.String()
	}
}

func FindInArr(str string, arr []string) bool {
	for _, s := range arr {
		if str == s {
			return true
		}
	}

	return false
}

func AddPostfix(str string, width int, ch string) string {
	lent := runewidth.StringWidth(str)

	if width > lent {
		postFix := strings.Repeat(ch, width-lent)
		str += postFix
	}

	return str
}

func GetNumbWidth(numb int) int {
	str := strconv.Itoa(numb)
	width := len(str)

	return width
}

func FormatStr(format string, val interface{}) (string, bool) {
	str := fmt.Sprintf(format, val)

	if strings.Index(str,"%!") == 0 {
		return "", false
	}

	return str, true
}

func InArray(need interface{}, arr []string) bool {
	for _,v := range arr{
		if need == v{
			return true
		}
	}
	return false
}

func AddPad(str string) string {
	if vari.Length > 0 {
		gap := vari.Length - len(str)
		if vari.LeftPad != "" {
			vari.LeftPad = vari.LeftPad[:1]
			pads := strings.Repeat(vari.LeftPad, gap)
			str = pads + str
		} else if vari.RightPad != "" {
			vari.RightPad = vari.RightPad[:1]
			pads := strings.Repeat(vari.RightPad, gap)
			str = str + pads
		} else {
			vari.LeftPad = " "
			pads := strings.Repeat(vari.LeftPad, gap)
			str = pads + str
		}
	}

	return str
}