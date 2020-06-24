package parser

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"reflect"
	"strconv"
	"time"
)

type Parser struct {
	Reader  *csv.Reader
	header  []string
	options Options
}

type Options struct {
	UseHeader  bool
	TimeLayout string
}

func New(inputStream io.Reader, options Options) Parser {
	csvReader := csv.NewReader(inputStream)

	var header []string
	var err error
	if options.UseHeader {
		header, err = csvReader.Read()
		if err != nil {
			log.Fatal(err)
		}
	}

	//csvReader.Comma = options.Comma
	//csvReader.Comment = options.Comment
	//csvReader.FieldsPerRecord = options.FieldsPerRecord
	//csvReader.LazyQuotes = options.LazyQuotes
	//csvReader.TrimLeadingSpace = options.TrimLeadingSpace
	//csvReader.ReuseRecord = options.ReuseRecord

	return Parser{
		Reader:  csvReader,
		header:  header,
		options: options,
	}
}

var boolValues = map[string]rune{
	"1":    1,
	"true": 1,
}

func (p *Parser) ReadInto(value interface{}) error {
	csvLine, err := p.reader.Read()
	if err != nil {
		return err
	}

	csvMap := toMap(p.header, csvLine)

	csvTagsMap := getCsvTags(value)

	for key, v := range csvMap {
		valueField, found := csvTagsMap[key]
		if !found {
			continue
		}

		valueFieldType := valueField.Type()

		switch valueField.Kind() {
		case reflect.Int:
			valueInt, err := strconv.Atoi(v)
			if err != nil {
				return errors.New(fmt.Sprintf("cannot parse %s into type %s", v, valueFieldType))
			}
			valueField.SetInt(int64(valueInt))
			break
		case reflect.String:
			valueField.SetString(v)
			break
		case reflect.Bool:
			_, valueBool := boolValues[v]
			valueField.SetBool(valueBool)
			break
		case reflect.Struct:
			switch valueFieldType.String() {
			case "time.Time":
				layout := p.options.TimeLayout
				t, err := time.Parse(layout, v)
				if err != nil {
					return errors.New(fmt.Sprintf("cannot parse %s into type %s", v, valueFieldType))
				}
				valueField.Set(reflect.ValueOf(t))
				break
			default:
				return errors.New(fmt.Sprintf("this lib cannot parse %s", valueFieldType))
			}
			break
		default:
			return errors.New(fmt.Sprintf("this lib cannot parse %s", valueFieldType))
		}
	}

	return nil
}

func getCsvTags(value interface{}) map[string]reflect.Value {
	valueReflect := reflect.ValueOf(value).Elem()
	valueTypeReflect := reflect.TypeOf(value).Elem()

	var csvTagsMap = map[string]reflect.Value{}

	for i := 0; i < valueTypeReflect.NumField(); i++ {
		field := valueTypeReflect.Field(i)

		tag := field.Tag.Get("csv")
		fieldName := field.Name

		valueField := valueReflect.FieldByName(fieldName)
		csvTagsMap[tag] = valueField
	}
	return csvTagsMap
}

func toMap(header []string, line []string) map[string]string {
	var out = make(map[string]string)

	for index, h := range header {
		out[h] = line[index]
	}

	return out
}
