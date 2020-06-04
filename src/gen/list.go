package gen

import (
	"fmt"
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	"regexp"
	"strconv"
	"strings"
)

func GenerateList(field *model.DefField) model.FieldValue {
	fieldValue := model.FieldValue{}
	GenerateListField(field, &fieldValue)

	return fieldValue
}

func GenerateListField(field *model.DefField, fieldValue *model.FieldValue) {
	fieldValue.Field = field.Field
	fieldValue.Precision = field.Precision

	if len(field.Fields) > 0 {
		for _, child := range field.Fields {
			childValue := model.FieldValue{}
			GenerateListField(&child, &childValue)
		}
	} else {
		GenerateFieldValues(field, fieldValue)
	}
}

func GenerateFieldValues(field *model.DefField, fieldValue *model.FieldValue) {
	if strings.Index(field.Range, ".txt") > -1 {
		GenerateFieldValuesFromText(field, fieldValue)
	} else {
		GenerateFieldValuesFromList(field, fieldValue)
	}
}

func GenerateFieldValuesFromList(field *model.DefField, fieldValue *model.FieldValue) {
	//rang := strings.TrimSpace(field.Range)
	rang := field.Range

	rangeItems := ParseRange(rang)

	index := 0
	for _, rangeItem := range rangeItems {
		if index >= constant.MaxNumb { break }
		if rangeItem == "" { continue }

		entry, stepStr, limit := ParseRangeItem(rangeItem)
		typ, desc := ParseEntry(entry)

		items := make([]interface{}, 0)
		if typ == "literal" {

		} else if typ == "interval" {
			elemArr := strings.Split(desc, "-")
			startStr := elemArr[0]
			endStr := startStr
			if len(elemArr) > 1 { endStr = elemArr[1] }
			items = GenerateValuesFromInterval(field, startStr, endStr, stepStr, limit)
		}

		fieldValue.Values = append(fieldValue.Values, items...)
		index = index + len(items)
	}

	if len(fieldValue.Values) == 0 {
		fieldValue.Values = append(fieldValue.Values, "N/A")
	}
}

func CheckRangeType(startStr string, endStr string, stepStr string) (string, interface{}, int, bool) {
	rand := false

	_, errInt1 := strconv.ParseInt(startStr, 0, 64)
	_, errInt2 := strconv.ParseInt(endStr, 0, 64)
	if errInt1 == nil && errInt2 == nil {
		var step interface{} = 1
		if strings.ToLower(strings.TrimSpace(stepStr)) != "r" {
			stepInt, errInt3 := strconv.Atoi(stepStr)
			if errInt3 == nil {
				step = stepInt
			}
		} else {
			rand = true
		}

		return "int", step, 0, rand

	} else {
		startFloat, errFloat1 := strconv.ParseFloat(startStr, 64)
		_, errFloat2 := strconv.ParseFloat(endStr, 64)
		if errFloat1 == nil && errFloat2 == nil {
			var step interface{} = 0.1
			if strings.ToLower(strings.TrimSpace(stepStr)) != "r" {
				stepFloat, errFloat3 := strconv.ParseFloat(stepStr, 64)
				if errFloat3 == nil {
					step = stepFloat
				}
			} else {
				rand = true
			}

			precision := getPrecision(startFloat, step)

			return "float", step, precision, rand

		} else if len(startStr) == 1 && len(endStr) == 1 {
			var step interface{} = 1
			if strings.ToLower(strings.TrimSpace(stepStr)) != "r" {
				stepChar, errChar3 := strconv.Atoi(stepStr)
				if errChar3 == nil {
					step = stepChar
				}
			} else {
				rand = true
			}

			return "char", step, 0, rand
		}
	}

	return "string", 1, 0, false
}

func ParseRange(rang string) []string {
	items := make([]string, 0)

	tagOpen := false
	temp := ""
	for i := 0; i < len(rang); i++ {
		c := rang[i]

		if int32(c) == constant.RightChar {
			tagOpen = false
		} else if int32(c) == constant.LeftChar  {
			tagOpen = true
		}

		if i == len(rang) - 1 {
			temp += fmt.Sprintf("%c", c)
			items = append(items, temp)
		} else if !tagOpen && c == ',' {
			items = append(items, temp)
			temp = ""
			tagOpen = false
		} else {
			temp += fmt.Sprintf("%c", c)
		}
	}

	return items
}

func ParseRangeItem(item string) (string, string, int) {
	entry := ""
	step := "1"
	limit := -1

	regx := regexp.MustCompile(`\{(.*)\}`)
	arr := regx.FindStringSubmatch(item)
	if len(arr) == 2 {
		limit, _ = strconv.Atoi(arr[1])
	}
	item = regx.ReplaceAllString(item, "")

	sectionArr := strings.Split(item, ":")
	entry = sectionArr[0]
	if len(sectionArr) == 2 {
		step = sectionArr[1]
	}

	return entry, step, limit
}

func GenerateValuesFromInterval(field *model.DefField, startStr string, endStr string, stepStr string, limit int) []interface{} {
	items := make([]interface{}, 0)

	dataType, step, precision, rand := CheckRangeType(startStr, endStr, stepStr)

	if dataType == "int" {
		startInt, _ := strconv.ParseInt(startStr, 0, 64)
		endInt, _ := strconv.ParseInt(endStr, 0, 64)

		items = GenerateIntItems(startInt, endInt, step, rand, limit)
	} else if dataType == "float" {
		startFloat, _ := strconv.ParseFloat(startStr, 64)
		endFloat, _ := strconv.ParseFloat(endStr, 64)
		field.Precision = precision

		items = GenerateFloatItems(startFloat, endFloat, step.(float64), rand, limit)
	} else if dataType == "char" {
		items = GenerateByteItems(byte(startStr[0]), byte(endStr[0]), step, rand, limit)
	} else if dataType == "string" {
		items = append(items, startStr)
		if startStr != endStr {
			items = append(items, endStr)
		}
	}

	return items
}

func ParseEntry(str string) (string, string) {
	typ := ""
	desc := ""

	str = strings.TrimSpace(str)
	if int32(str[0]) == constant.LeftChar {
		typ = "literal"
		desc = strings.ReplaceAll(str, string(constant.LeftChar), "")
		desc = strings.ReplaceAll(desc,string(constant.RightChar), "")
	} else {
		typ = "interval"
		desc = str
	}

	return typ, desc
}