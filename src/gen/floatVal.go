package gen

import (
	commonUtils "github.com/easysoft/zendata/src/utils/common"
	constant "github.com/easysoft/zendata/src/utils/const"
	"strconv"
	"strings"
)

func GenerateFloatItems(start float64, end float64, step interface{}, rand bool, repeat int) []interface{} {
	if !rand {
		return GenerateFloatItemsByStep(start, end, step.(float64), repeat)
	} else{
		return GenerateFloatItemsRand(start, end, step.(float64), repeat)
	}
}

func GenerateFloatItemsByStep(start float64, end float64, step float64, repeat int) []interface{} {
	arr := make([]interface{}, 0)

	total := 0

	for i := 0; true; {
		val := start + float64(i) * step
		if (val > end && step > 0) || (val < end && step < 0)  {
			break
		}

		for round := 0; round < repeat; round++ {
			arr = append(arr, val)

			total++
			if total > constant.MaxNumb {
				break
			}
		}
		if total > constant.MaxNumb {
			break
		}
		i++
	}

	return arr
}

func GenerateFloatItemsRand(start float64, end float64, step float64, repeat int) []interface{} {
	arr := make([]interface{}, 0)

	countInRound := (end - start) / step
	total := 0
	for i := float64(0); i < countInRound; {
		rand := commonUtils.RandNum64(int64(countInRound))
		if step < 0 {
			rand = rand * -1
		}

		val := start + float64(rand) * step

		for round := 0; round < repeat; round++ {
			arr = append(arr, val)

			total++
			if total > constant.MaxNumb {
				break
			}
		}

		if total > constant.MaxNumb {
			break
		}
		i++
	}

	return arr
}

func GetPrecision(base float64, step interface{}) (precision int, newStep float64) {
	val := base

	if step != nil {
		val += step.(float64)
	}

	str1 := strconv.FormatFloat(base, 'f', -1, 64)
	str2 := strconv.FormatFloat(val, 'f', -1, 64)

	index1 := strings.LastIndex(str1, ".")
	index2 := strings.LastIndex(str2, ".")

	if index1 < index2 {
		precision = len(str1) - index1 - 1
	} else {
		precision = len(str2) - index2 - 1
	}

	if step == nil || step == 0 {
		newStep = float64(1)
		for i := 0; i < precision; i++ {
			newStep = newStep / 10
		}
	} else {
		newStep = step.(float64)
	}

	return
}

func InterfaceToStr(val interface{}) string {
	str := "n/a"

	switch val.(type) {
		case int64:
			str = strconv.FormatInt(val.(int64), 10)
		case float64:
			precision, _ := GetPrecision(val.(float64), nil)
			str = strconv.FormatFloat(val.(float64), 'f', precision, 64)
		case byte:
			str = string(val.(byte))
		case string:
			str = val.(string)
		default:
	}
	return str
}