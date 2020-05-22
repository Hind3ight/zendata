package action

import (
	"github.com/easysoft/zendata/src/model"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	"github.com/easysoft/zendata/src/utils/vari"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

func ParseSql(file string, out string) {
	startTime := time.Now().Unix()
	vari.InputDir = filepath.Dir(file) + string(os.PathSeparator)

	statements := getCreateStatement(file)
	for tableName, statement := range statements {
		columns := getColumnsFromCreateStatement(statement)

		def := model.DefSimple{}
		def.Init(tableName, "automated export", "", "1.0")

		for _, col := range columns {
			field := model.FieldSimple{}
			field.Init(col)
			def.Fields = append(def.Fields, field)
		}

		bytes, _ := yaml.Marshal(def)
		out = fileUtils.UpdateDir(out)
		outFile := out + tableName + ".yaml"
		WriteToFile(outFile, string(bytes))
	}

	entTime := time.Now().Unix()
	logUtils.Screen(i118Utils.I118Prt.Sprintf("generate_yaml", len(statements), out, entTime - startTime ))
}

func getCreateStatement(file string) map[string]string {
	statements := map[string]string{}

	content, err := ioutil.ReadFile(file)
	if err != nil {
		logUtils.Screen(i118Utils.I118Prt.Sprintf("fail_to_read_file", file))
		return statements
	}

	re := regexp.MustCompile("(?siU)(CREATE TABLE.*;)")
	arr := re.FindAllString(string(content), -1)
	for _, item := range arr {
		re := regexp.MustCompile("(?iU)CREATE TABLE.*`(.+)`")
		arr2 := re.FindAllStringSubmatch(item, -1)

		statements[arr2[0][1]] = item
	}

	return statements
}

func getColumnsFromCreateStatement(sent string) []string {
	fieldLines := make([]string, 0)

	re := regexp.MustCompile("(?iU)`(.+)`\\s.*,")
	arr := re.FindAllStringSubmatch(string(sent), -1)
	for _, item := range arr {
		fieldLines = append(fieldLines, item[1])
	}

	return fieldLines
}