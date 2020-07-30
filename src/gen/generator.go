package gen

import (
	"fmt"
	"github.com/easysoft/zendata/src/model"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"github.com/easysoft/zendata/src/utils/vari"
	"github.com/mattn/go-runewidth"
	"strconv"
	"strings"
)

func GenerateForDefinition(defaultFile, configFile string, fieldsToExport *[]string, total int) ([][]string, []bool) {
	vari.DefaultDir = fileUtils.GetAbsDir(defaultFile)
	vari.ConfigDir = fileUtils.GetAbsDir(configFile)

	vari.Def = LoadConfigDef(defaultFile, configFile, fieldsToExport)
	vari.Res = LoadResDef(*fieldsToExport)

	topFieldNameToValuesMap := map[string][]string{}
	colIsNumArr := make([]bool, 0)

	// 为每个field生成值列表
	for index, field := range vari.Def.Fields {
		if !stringUtils.FindInArr(field.Field, *fieldsToExport) {
			continue
		}

		values := GenerateForField(&field, total, true)
		vari.Def.Fields[index].Precision = field.Precision

		topFieldNameToValuesMap[field.Field] = values
		colIsNumArr = append(colIsNumArr, field.IsNumb)
	}

	// 处理数据
	arrOfArr := make([][]string, 0) // 2 dimension arr for child, [ [a,b,c], [1,2,3] ]
	for _, child := range vari.Def.Fields {
		if !stringUtils.FindInArr(child.Field, *fieldsToExport) {
			continue
		}

		childValues := topFieldNameToValuesMap[child.Field]
		arrOfArr = append(arrOfArr, childValues)
	}

	rows := putChildrenToArr(arrOfArr, total)
	return rows, colIsNumArr
}

func GenerateForField(field *model.DefField, total int, withFix bool) (values []string) {
	if len(field.Fields) > 0 { // sub fields
		arrOfArr := make([][]string, 0) // 2 dimension arr for child, [ [a,b,c], [1,2,3] ]
		for _, child := range field.Fields {
			childValues := GenerateForField(&child, total, withFix)
			arrOfArr = append(arrOfArr, childValues)
		}

		count := total
		count = getRecordCount(arrOfArr)
		if count > total {
			count = total
		}
		values = connectChildrenToSingleStr(arrOfArr, count)
		values = LoopSubFields(field, values, count, true)

	} else if field.From != "" { // refer to res

		if field.Use != "" { // refer to instance
			groupValues := vari.Res[field.From]
			groups := strings.Split(field.Use, ",")
			for _, group := range groups {
				if group == "all" {
					for _, arr := range groupValues { // add all
						values = append(values, arr...)
					}
				} else {
					values = append(values, groupValues[group]...)
				}
			}
		} else if field.Select != "" { // refer to excel
			groupValues := vari.Res[field.From]
			slct := field.Select
			values = append(values, groupValues[slct]...)
		}

		values = LoopSubFields(field, values, total, true)

	} else if field.Config != "" { // refer to another define
		groupValues := vari.Res[field.Config]
		values = append(values, groupValues["all"]...)

	} else { // leaf field
		values = GenerateFieldItemsFromDefinition(field)
	}

	return values
}

func GenerateFieldItemsFromDefinition(field *model.DefField) []string {
	values := make([]string, 0)

	fieldWithValues := GenerateList(field)

	index := 0
	count := 0
	for {
		// 处理格式、前后缀、loop等
		val := GenerateFieldValWithFix(field, fieldWithValues, &index, true)
		values = append(values, val)

		count++
		if index >= len(fieldWithValues.Values) || count >= vari.Total {
			break
		}
	}

	return values
}

func GetFieldValStr(field model.DefField, val interface{}) string {
	str := "n/a"
	success := false

	switch val.(type) {
		case int64:
			if field.Format != "" {
				str, success = stringUtils.FormatStr(field.Format, val.(int64))
			}
			if !success {
				str = strconv.FormatInt(val.(int64), 10)
			}
		case float64:
			precision := 0
			if field.Precision > 0 {
				precision = field.Precision
			}
			if field.Format != "" {
				str, success = stringUtils.FormatStr(field.Format, val.(float64))
			}
			if !success {
				str = strconv.FormatFloat(val.(float64), 'f', precision, 64)
			}
		case byte:
			str = string(val.(byte))
			if field.Format != "" {
				str, success = stringUtils.FormatStr(field.Format, str)
			}
			if !success {
				str = string(val.(byte))
			}
		case string:
			str = val.(string)
			fmt.Sprintf(str)
		default:
	}

	return str
}

func LoopSubFields(field *model.DefField, oldValues []string, total int, withFix bool) (values []string) {

	fieldValue := model.FieldWithValues{}

	for _, val := range oldValues {
		fieldValue.Values = append(fieldValue.Values, val)
	}

	indexOfRow := 0
	count := 0
	for {
		// 处理格式、前后缀、loop等
		str := GenerateFieldValWithFix(field, fieldValue, &indexOfRow, withFix)
		values = append(values, str)

		count++
		if indexOfRow >= len(fieldValue.Values) || count >= total {
			break
		}
	}

	return
}

func GenerateFieldValWithFix(field *model.DefField, fieldValue model.FieldWithValues,
		indexOfRow *int, withFix bool) (loopStr string) {
	prefix := field.Prefix
	postfix := field.Postfix

	computerLoop(field)

	for j := 0; j < (*field).LoopIndex; j++ {
		if loopStr != "" {
			loopStr = loopStr + field.Loopfix
		}

		str := GenerateFieldVal(*field, fieldValue, indexOfRow)
		loopStr = loopStr + str

		*indexOfRow++
	}

	if withFix && !vari.Trim {
		loopStr = prefix + loopStr + postfix
	}

	if field.Width > runewidth.StringWidth(loopStr) {
		loopStr = stringUtils.AddPad(loopStr, *field)
	}

	(*field).LoopIndex = (*field).LoopIndex + 1
	if (*field).LoopIndex > (*field).LoopEnd {
		(*field).LoopIndex = (*field).LoopStart
	}

	return
}

func GenerateFieldVal(field model.DefField, fieldValue model.FieldWithValues, index *int) (val string) {
	// 叶节点
	idx := *index % len(fieldValue.Values)
	str := fieldValue.Values[idx]
	val = GetFieldValStr(field, str)

	return
}

func computerLoop(field *model.DefField) {
	if (*field).LoopIndex != 0 {
		return
	}

	arr := strings.Split(field.Loop, "-")
	(*field).LoopStart, _ = strconv.Atoi(arr[0])
	if len(arr) > 1 {
		field.LoopEnd, _ = strconv.Atoi(arr[1])
	}

	if (*field).LoopStart == 0 {
		(*field).LoopStart = 1
	}
	if (*field).LoopEnd == 0 {
		(*field).LoopEnd = 1
	}

	(*field).LoopIndex = (*field).LoopStart
}

func putChildrenToArr(arrOfArr [][]string, total int) (values [][]string) {
	indexArr := make([]int, 0)
	if vari.Recursive {
		indexArr = getModArr(arrOfArr)
	}

	for i := 0; i < total; i++ {
		strArr := make([]string, 0)
		for j := 0; j < len(arrOfArr); j++ {
			child := arrOfArr[j]

			var index int
			if vari.Recursive {
				mod := indexArr[j]
				index = i / mod % len(child)
			} else {
				index = i % len(child)
			}

			strArr = append(strArr, child[index])
		}

		values = append(values, strArr)
	}

	return
}

func connectChildrenToSingleStr(arrOfArr [][]string, total int) (ret []string)  {
	valueArr := putChildrenToArr(arrOfArr, total)

	for _, arr := range valueArr {
		ret = append(ret, strings.Join(arr, ""))
	}
	return
}

func getRecordCount(arrOfArr [][]string) int {
	count := 1
	for i := 0; i < len(arrOfArr); i++ {
		arr := arrOfArr[i]
		count = len(arr) * count
	}
	return count
}

func getModArr(arrOfArr [][]string) []int {
	indexArr := make([]int, 0)
	for _, _ = range arrOfArr {
		indexArr = append(indexArr, 0)
	}

	for i := 0; i < len(arrOfArr); i++ {
		loop := 1
		for j := i + 1; j < len(arrOfArr); j++ {
			loop = loop * len(arrOfArr[j])
		}

		indexArr[i] = loop
	}

	return indexArr
}
