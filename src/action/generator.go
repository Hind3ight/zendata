package action

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/easysoft/zendata/src/gen"
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"github.com/easysoft/zendata/src/utils/vari"
	"github.com/fatih/color"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	FileWriter *os.File
	HttpWriter http.ResponseWriter
)

func Generate(defaultFile string, configFile string, total int, fieldsToExportStr string, out string, format string, table string) {
	startTime := time.Now().Unix()

	if defaultFile != "" && configFile == "" {
		configFile = defaultFile
		defaultFile = ""
	}

	fieldsToExport := make([]string, 0)
	if fieldsToExportStr != "" {
		fieldsToExport = strings.Split(fieldsToExportStr, ",")
	}

	vari.Total = total

	rows, colIsNumArr, err := gen.GenerateForDefinition(defaultFile, configFile, &fieldsToExport, total)
	if err != nil {
		return
	}
	Print(rows, format, table, colIsNumArr, fieldsToExport)

	entTime := time.Now().Unix()
	if vari.RunMode == constant.RunModeServerRequest {
		logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("server_response", len(rows), entTime - startTime), color.FgCyan)
	}
}

func Print(rows [][]string, format string, table string, colIsNumArr []bool, fields []string) {
	if format == constant.FormatText {
		printTextHeader(fields)
	} else if format == constant.FormatSql {
		printSqlHeader(fields, table)
	} else if format == constant.FormatJson {
		printJsonHeader()
	} else if format == constant.FormatXml {
		printXmlHeader(fields, table)
	}

	for i, cols := range rows {
		row := make([]string, 0)
		rowXml := model.XmlRow{}
		valuesForSql := make([]string, 0)
		lineForText := ""

		for j, col := range cols {
			col = replacePlaceholder(col)

			lineForText = lineForText + col

			row = append(row, col)
			rowXml.Cols = append(rowXml.Cols, col)

			colVal := stringUtils.ConvertForSql(col)
			if !colIsNumArr[j] { colVal = "'" + colVal + "'" }
			valuesForSql = append(valuesForSql, colVal)
		}

		if format == constant.FormatText {
			printLine(lineForText)
		} else if format == constant.FormatSql {
			printLine(genSqlLine(strings.Join(valuesForSql, ", "), i, len(rows)))
		} else if format == constant.FormatJson {
			printLine(genJsonLine(i, row, len(rows), fields))
		} else if format == constant.FormatXml {
			printLine(getXmlLine(i, rowXml, len(rows)))
		}
	}
}

func printTextHeader(fields []string) {
	if !vari.WithHead {
		return
	}
	headerLine := ""
	for idx, field := range fields {
		headerLine += field
		if idx < len(fields)-1 {
			headerLine += vari.HeadSep
		}
	}

	printLine(headerLine)
}

func printSqlHeader(fields []string, table string) {
	fieldNames := make([]string, 0)
	for _, f := range fields { fieldNames = append(fieldNames, "`" + f + "`") }
	printLine(fmt.Sprintf("INSERT INTO %s(%s)", table, strings.Join(fieldNames, ", ")))
}

func printJsonHeader() {
	printLine("[")
}

func printXmlHeader(fields []string, table string) {
	printLine("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<testdata>\n  <title>Test Data</title>")
}

func RowToJson(cols []string, fieldsToExport []string) string {
	rowMap := map[string]string{}
	for j, col := range cols {
		rowMap[fieldsToExport[j]] = col
	}

	jsonObj, _ := json.Marshal(rowMap)
	respJson := string(jsonObj)

	return respJson
}

func genSqlLine(valuesForSql string, i int, length int) string {

	temp := ""
	if i == 0 {
		temp = fmt.Sprintf("  VALUES (%s)", valuesForSql)
	} else {
		temp = fmt.Sprintf("         (%s)", valuesForSql)
	}

	if i < length - 1 {
		temp = temp + ", "
	} else {
		temp = temp + "; "
	}

	return temp
}

func genJsonLine(i int, row []string,  length int, fields []string) string {
	temp := "  " + RowToJson(row, fields)
	if i < length - 1 {
		temp = temp + ", "
	} else {
		temp = temp + "\n]"
	}

	return temp
}

func getXmlLine(i int, rowXml model.XmlRow, length int) string {
	rowXmlStr, _ := xml.Marshal(rowXml)
	text := "  " + string(rowXmlStr)
	if i == length - 1 {
		text = text + "\n</testdata>"
	}
	return text
}

func replacePlaceholder(col string) string {
	ret := col

	re := regexp.MustCompile("(?siU)\\${(.*)}")
	arr := re.FindAllStringSubmatch(col, -1)

	strForReplaceMap := map[string][]string{}
	for _, childArr := range arr {
		placeholderStr := childArr[1]
		strForReplaceMap[placeholderStr] = getValForPlaceholder(placeholderStr, len(childArr))

		for _, str := range strForReplaceMap[placeholderStr] {
			temp := gen.Placeholder(placeholderStr)
			ret = strings.Replace(ret, temp, str, 1)
		}
	}

	return ret
}

func getValForPlaceholder(placeholderStr string, count int) []string {
	mp := vari.RandFieldNameToValuesMap[placeholderStr]
	tp := mp["type"].(string)
	repeatObj := mp["repeat"]
	repeat := "1"
	if repeatObj != nil {
		repeat = repeatObj.(string)
	}

	strs := make([]string, 0)
	if tp == "int" {
		start := mp["start"].(string)
		end := mp["end"].(string)
		precision := mp["precision"].(string)

		strs = gen.GetRandFromRange("int", start, end, "1", repeat, precision, count)
	} else if tp == "float" {
		start := mp["start"].(string)
		end := mp["end"].(string)
		precision := mp["precision"].(string)

		strs = gen.GetRandFromRange("float", start, end, "1", repeat, precision, count)
	} else if tp == "char" {
		start := mp["start"].(string)
		end := mp["end"].(string)
		precision := mp["precision"].(string)

		strs = gen.GetRandFromRange("char", start, end, "1", repeat, precision, count)
	} else if tp == "list" {
		list := mp["list"].([]string)
		strs = gen.GetRandFromList(list, repeat, count)
	}

	return strs
}

func PrintErrMsg(msg string) {
	logUtils.PrintToWithColor(msg, color.FgCyan)
}

func printLine(line string) {
	if FileWriter != nil {
		PrintToFile(line)
	} else if vari.RunMode == constant.RunModeServerRequest {
		PrintToHttp(line)
	} else {
		PrintToScreen(line)
	}
}
func PrintToFile(line string) {
	fmt.Fprintln(FileWriter, line)
}
func PrintToHttp(line string) {
	fmt.Fprintln(HttpWriter, line)
}
func PrintToScreen(line string) {
	fmt.Println(line)
}