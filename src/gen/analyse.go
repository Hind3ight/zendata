package gen

import (
	"encoding/json"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	"github.com/easysoft/zendata/src/utils/vari"
	"strings"
)

const (

)

func Analyse(output, fieldsToExportStr, configFile, defaultFile string) {


	fieldsToExport := make([]string, 0)
	if fieldsToExportStr != "" {
		fieldsToExport = strings.Split(fieldsToExportStr, ",")
	}

	vari.Def = LoadConfigDef(defaultFile, configFile, &fieldsToExport)
	vari.Res = LoadResDef(fieldsToExport)

	data := fileUtils.ReadFile(output)

	mapArr := LinesToMap(data, fieldsToExport)
	jsonObj, _ := json.Marshal(mapArr)
	respJson := string(jsonObj)

	path := output[:strings.LastIndex(output, ".")] + ".json"
	fileUtils.WriteFile(path, respJson)
}

func LinesToMap(str string, fieldsToExport []string) (ret []map[string]string) {
	start := 0
	if vari.WithHead {
		start = 1
	}

	ret = []map[string]string{}

	for index, line := range strings.Split(str, "\n") {
		if index < start {
			continue
		}

		rowMap := map[string]string{}
		left := line

		for j, field := range vari.Def.Fields {
			col := left[:field.Length]
			rowMap[fieldsToExport[j]] = col

			left = left[field.Length:]
		}

		ret = append(ret, rowMap)
	}
	return
}