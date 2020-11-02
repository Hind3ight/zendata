package gen

import (
	"encoding/csv"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	constant "github.com/easysoft/zendata/src/utils/const"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"github.com/easysoft/zendata/src/utils/vari"
	"github.com/mattn/go-runewidth"
)

const (
	sheetName = "Sheet1"
)
var (
	csvWriter *csv.Writer
)

func Write(rows [][]string, format string, table string, colIsNumArr []bool,
		fields []string) (lines []interface{}) {

	f := excelize.NewFile()
	index := f.NewSheet(sheetName)
	f.SetActiveSheet(index)

	if format == constant.FormatExcel {
		printExcelHeader(fields, f)
	} else if format == constant.FormatCsv {
		csvWriter = csv.NewWriter(logUtils.FileWriter)
	}

	csvData := make([][]string, 0)
	for i, cols := range rows {
		csvRow := make([]string, 0)

		for j, col := range cols {
			col = replacePlaceholder(col)
			field := vari.TopFieldMap[fields[j]]
			if field.Length > runewidth.StringWidth(col) {
				col = stringUtils.AddPad(col, field)
			}

			if format == constant.FormatExcel {
				colName, _ := excelize.CoordinatesToCellName(j + 1, i + 2)
				f.SetCellValue(sheetName, colName, col)

			} else if format == constant.FormatCsv {
				csvRow = append(csvRow, col)
			}
		}
		csvData = append(csvData, csvRow)
	}

	var err error
	if format == constant.FormatExcel {
		err = f.SaveAs(logUtils.FilePath)
	} else if format == constant.FormatCsv {
		err = csvWriter.WriteAll(csvData)
		csvWriter.Flush()
	}
	if err != nil {
		logUtils.PrintErrMsg(err.Error())
	}

	return
}

func printExcelHeader(fields []string, f *excelize.File) {
	headerLine := ""
	for idx, field := range fields {
		colName, _ := excelize.CoordinatesToCellName(idx + 1, 1)
		f.SetCellValue(sheetName, colName, field)
	}

	logUtils.PrintLine(headerLine)
}