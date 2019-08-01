package test

import (
	"clearance-adapter/factory"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	factory.Init()
	initCSL()
	initP2BrandDB()
}

func initP2BrandDB() {
	initStore()
	initProduct()
	initLocation()
}

func readDataFromCSV(filename string) (map[int]string, [][]string) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err.Error())
	}
	reader := csv.NewReader(strings.NewReader(string(bytes)))
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err.Error())
	}

	headerRow := records[0]
	headers := make(map[int]string, 0)
	for i, header := range headerRow {
		headers[i] = header
	}

	return headers, records[1:len(records)]
}

func setObjectValue(headers map[int]string, data []string, obj interface{}) {
	for colIndex, val := range data {
		header := headers[colIndex]
		obj := reflect.ValueOf(obj)
		if obj.Kind() == reflect.Ptr && !obj.Elem().CanSet() {
			continue
		}
		field := obj.Elem().FieldByName(header)
		if !field.IsValid() {
			continue
		}
		if field.Kind() == reflect.String {
			field.SetString(val)
		} else if field.Kind() == reflect.Int64 {
			v, _ := strconv.ParseInt(val, 10, 64)
			field.SetInt(v)
		} else if field.Kind() == reflect.Int {
			v, _ := strconv.ParseInt(val, 10, 64)
			field.SetInt(v)
		} else if field.Kind() == reflect.Bool {
			if val == "1" {
				field.SetBool(true)
			} else {

				field.SetBool(false)
			}
		} else if field.Kind() == reflect.Float64 {
			v, _ := strconv.ParseFloat(val, 64)
			field.SetFloat(v)
		}
	}
}
