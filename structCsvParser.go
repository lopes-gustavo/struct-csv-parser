package parser

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
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

func New(inputStream io.Reader, options Options) (Parser, error) {
	csvReader := csv.NewReader(inputStream)

	var header []string
	var err error
	if options.UseHeader {
		header, err = csvReader.Read()
		if err != nil {
			return Parser{}, err
		}
	}

	return Parser{
		Reader:  csvReader,
		header:  header,
		options: options,
	}, nil
}

var boolValues = map[string]rune{
	"1":    1,
	"true": 1,
}

func (p *Parser) ReadInto(value interface{}) error {
	csvValuesSlice, err := p.Reader.Read()
	if err != nil {
		return err
	}

	csvValuesMap := p.toMap(csvValuesSlice)
	csvFieldNameToCsvTagMap := getCsvTags(value)

	for csvHeader, csvValue := range csvValuesMap {

		valueReflect := reflect.ValueOf(value).Elem()

		var valueField reflect.Value

		if p.options.UseHeader {
			fieldName, found := csvFieldNameToCsvTagMap[csvHeader]
			if !found {
				continue
			}
			valueField = valueReflect.FieldByName(fieldName)
		} else {
			fieldNum, _ := strconv.Atoi(csvHeader)
			valueField = valueReflect.Field(fieldNum)
		}

		valueFieldType := valueField.Type()

		switch valueField.Kind() {
		case reflect.Int:
			valueInt, err := strconv.Atoi(csvValue)
			if err != nil {
				return errors.New(fmt.Sprintf("cannot parse %s into type %s", csvValue, valueFieldType))
			}
			valueField.SetInt(int64(valueInt))
			break
		case reflect.String:
			valueField.SetString(csvValue)
			break
		case reflect.Bool:
			_, valueBool := boolValues[csvValue]
			valueField.SetBool(valueBool)
			break
		case reflect.Struct:
			switch valueFieldType.String() {
			case "time.Time":
				layout := p.options.TimeLayout
				t, err := time.Parse(layout, csvValue)
				if err != nil {
					return errors.New(fmt.Sprintf("cannot parse %s into type %s", csvValue, valueFieldType))
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

func getCsvTags(value interface{}) map[string]string {
	valueTypeReflect := reflect.TypeOf(value).Elem()

	var csvTagsMap = map[string]string{}

	for i := 0; i < valueTypeReflect.NumField(); i++ {
		field := valueTypeReflect.Field(i)

		tag := field.Tag.Get("csv")
		fieldName := field.Name

		csvTagsMap[tag] = fieldName
	}
	return csvTagsMap
}

func (p *Parser) toMap(line []string) map[string]string {
	var out = make(map[string]string)

	if p.options.UseHeader {
		for index, h := range p.header {
			out[h] = line[index]
		}
	} else {
		for index := range line {
			out[strconv.Itoa(index)] = line[index]
		}
	}

	return out
}
