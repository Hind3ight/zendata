package gen

import (
	"database/sql"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func GenerateFieldValuesFromExcel(path string, field *model.DefField) (map[string][]string, string) {
	values := map[string][]string{}

	idx := strings.LastIndex(field.From, ".")
	arr := strings.Split(field.From, ".")
	dbName := arr[len(arr) - 2]
	tableName := field.From[idx + 1:]

	list := make([]string, 0)
	selectCol := ""
	ConvertExcelToSQLiteIfNeeded(dbName, path)

	list, selectCol = ReadDataFromSQLite(*field, dbName, tableName)

	// get step and rand
	rand := false
	step := 1
	// get index for data retrieve
	numbs := GenerateIntItems(0, (int64)(len(list)-1), step, rand)
	// get data by index
	index := 0
	for _, numb := range numbs {
		item := list[numb.(int64)]

		if index >= constant.MaxNumb {
			break
		}
		if strings.TrimSpace(item) == "" {
			continue
		}

		values[selectCol] = append(values[selectCol], item)
		index = index + 1
	}

	return values, dbName
}

func ConvertExcelToSQLiteIfNeeded(dbName string, path string) {
	excel, err := excelize.OpenFile(path)
	if err != nil {
		logUtils.Screen(i118Utils.I118Prt.Sprintf("fail_to_read_file", path))
		return
	}

	if !isExcelChanged(path) {
		return
	}

	for _, sheet := range excel.GetSheetList() {
		rows, err := excel.GetRows(sheet)

		dropTemplate := `DROP TABLE IF EXISTS %s;`
		ddlTemplate := `CREATE TABLE %s (
						%s
					);`
		insertTemplate := "INSERT INTO %s (%s) VALUES %s"

		colDefine := ""
		colList := ""
		index := 0
		for _, col := range rows[0] {
			val := strings.TrimSpace(col)
			if index > 0 {
				colDefine = colDefine + ",\n"
				colList = colList + ", "
			}

			colProp := ""
			if val == "seq" {
				colProp = "CHAR (5) PRIMARY KEY ASC UNIQUE"
			} else {
				colProp = "VARCHAR"
			}
			colDefine = "    " + colDefine + val + " " + colProp

			colList = colList + val
			index++
		}

		valList := ""
		for rowIndex, row := range rows {
			if rowIndex == 0 {
				continue
			}

			if rowIndex > 1 {
				valList = valList + ", "
			}
			valList = valList + "("

			for colIndex, colCell := range row {
				if colIndex > 0 {
					valList = valList + ", "
				}
				valList = valList + "'" + colCell + "'"
			}
			valList = valList + ")"
		}

		tableName := dbName + "_" + sheet
		dropSql := fmt.Sprintf(dropTemplate, tableName)
		ddl := fmt.Sprintf(ddlTemplate, tableName, colDefine)
		insertSql := fmt.Sprintf(insertTemplate, tableName, colList, valList)

		db, err := sql.Open("sqlite3", constant.SqliteSource)
		defer db.Close()
		_, err = db.Exec(dropSql)
		_, err = db.Exec(ddl)
		if err != nil {
			logUtils.Screen(i118Utils.I118Prt.Sprintf("fail_to_create_table", tableName, err.Error()))
			return
		} else {
			_, err = db.Exec(insertSql)
			if err != nil {
				logUtils.Screen(i118Utils.I118Prt.Sprintf("fail_to_exec_query", insertSql, err.Error()))
				return
			}
		}

	}
}

func ReadDataFromSQLite(field model.DefField, dbName string, tableName string) ([]string, string) {
	list := make([]string, 0)

	db, err := sql.Open(constant.SqliteDriver, constant.SqliteSource)
	defer db.Close()
	if err != nil {
		logUtils.Screen(i118Utils.I118Prt.Sprintf("fail_to_connect_sqlite", constant.SqliteSource, err.Error()))
		return list, ""
	}

	selectCol := field.Select
	from := dbName + "_" + tableName
	where := field.Where
	if !strings.Contains(where, "LIMIT") {
		where = where + " LIMIT " + strconv.Itoa(constant.MaxNumb)
	}

	sqlStr := fmt.Sprintf("SELECT %s FROM %s WHERE %s", selectCol, from, where)
	rows, err := db.Query(sqlStr)
	if err != nil {
		logUtils.Screen(i118Utils.I118Prt.Sprintf("fail_to_exec_query", sqlStr, err.Error()))
		return list, ""
	}

	valMapArr := make([]map[string]string, 0)
	columns, err := rows.Columns()
	colNum := len(columns)

	colIndexToName := map[int]string{}
	for index, col := range columns {
		colIndexToName[index] = col
	}

	var values = make([]interface{}, colNum)
	for i, _ := range values {
		var itf string
		values[i] = &itf
	}

	for rows.Next() {
		err = rows.Scan(values...)
		if err != nil {
			logUtils.Screen(i118Utils.I118Prt.Sprintf("fail_to_parse_row", err.Error()))
			return list, ""
		}

		rowMap := map[string]string{}
		for index, v := range values {
			item := v.(*string)

			rowMap[colIndexToName[index]] = *item
		}

		valMapArr = append(valMapArr, rowMap)
	}

	for _, item := range valMapArr {
		idx := 0
		for _, val := range item {
			if idx > 0 { break }
			list = append(list, val)
			idx++
		}
	}

	return list, selectCol
}

func isExcelChanged(path string) bool {
	db, err := sql.Open(constant.SqliteDriver, constant.SqliteSource)
	defer db.Close()
	if err != nil {
		logUtils.Screen(i118Utils.I118Prt.Sprintf("fail_to_connect_sqlite", constant.SqliteSource, err.Error()))
		return true
	}

	sqlStr := "SELECT id, name, changeTime FROM " + constant.SqliteTrackTable
	rows, err := db.Query(sqlStr)
	if err != nil {
		logUtils.Screen(i118Utils.I118Prt.Sprintf("fail_to_exec_query", sqlStr, err.Error()))
		return true
	}

	fileChangeTime := getFileModTime(path).Unix()

	found := false
	changed := false
	for rows.Next() {
		found = true
		var id int
		var name string
		var changeTime int64

		err = rows.Scan(&id, &name, &changeTime)
		if err != nil {
			logUtils.Screen(i118Utils.I118Prt.Sprintf("fail_to_parse_row", err.Error()))
			changed = true
			break
		}

		if path == name && changeTime < fileChangeTime {
			changed = true
			break
		}
	}
	rows.Close()

	if !found { changed = true }

	if changed {
		if !found {
			sqlStr = fmt.Sprintf("INSERT INTO excel_change(name, changeTime) VALUES('%s', %d)", path, fileChangeTime)
		} else {
			sqlStr = fmt.Sprintf("UPDATE excel_change SET changeTime = %d WHERE name = '%s'", fileChangeTime, path)
		}

		_, err = db.Exec(sqlStr)
		if err != nil {
			logUtils.Screen(i118Utils.I118Prt.Sprintf("fail_to_exec_query", sqlStr, err.Error()))
		}
	}

	return changed
}

func getFileModTime(path string) time.Time {
	f, err := os.Open(path)
	if err != nil {
		log.Println("open file error:" + path)
		return time.Now()
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		log.Println("stat fileinfo error")
		return time.Now()
	}

	fileChangeTime := fi.ModTime()

	timeStr := fileChangeTime.Format("2006-01-02 15:04:05")
	logUtils.Screen(i118Utils.I118Prt.Sprintf("file_change_time", timeStr))

	return fileChangeTime
}