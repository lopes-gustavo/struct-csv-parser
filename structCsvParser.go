// Package structCsvParser parses a csv into a struct instead of the default slice
// It uses the builtin encoding/csv for primary parsing
package structCsvParser

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"time"
)

// Parser is the main struct. It should not be created directly, but with structCsvParser.New
type Parser struct {
	// Reader is the underlying csv.Reader
	// Should be edited to change csv reading properties, like Comma separator or Quotes
	Reader     *csv.Reader
	header     []string
	options    Options
	boolValues map[string]rune
}

type Options struct {
	UseHeader  bool
	TimeLayout string
	BoolValues []string
}

var defaultBoolValues = map[string]rune{
	"1":    1,
	"true": 1,
}

// sliceToMap is a helper function to create a map in which the keys are the strings passed to it
// Its purpose is indexing, for faster reading
func sliceToMap(ss []string) map[string]rune {
	var out = map[string]rune{}
	for _, s := range ss {
		out[s] = 1
	}
	return out
}

// Creates a new parser
func New(reader io.Reader, options Options) (Parser, error) {
	csvReader := csv.NewReader(reader)

	var header []string
	var err error
	if options.UseHeader {
		header, err = csvReader.Read()
		if err != nil {
			return Parser{}, err
		}
	}

	boolValues := defaultBoolValues
	if options.BoolValues != nil {
		boolValues = sliceToMap(options.BoolValues)
	}

	return Parser{
		Reader:     csvReader,
		header:     header,
		options:    options,
		boolValues: boolValues,
	}, nil
}

// Reads one line from the csv and tries to put it into the provided struct, which must be passed as reference
// value must be a pointer to a struct
func (p *Parser) ReadInto(value interface{}) error {
	csvValuesSlice, err := p.Reader.Read()
	if err != nil {
		return err
	}

	csvValuesMap := p.toMap(csvValuesSlice)
	csvFieldNameToCsvTagMap := getFieldNamesFromCsvTag(value)

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
			_, valueBool := defaultBoolValues[csvValue]
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

// reflects into the provided struct looking for the `csv` tag
// Ignores if the csv tag is not provided or is "-"
func getFieldNamesFromCsvTag(value interface{}) map[string]string {
	valueTypeReflect := reflect.TypeOf(value).Elem()

	var csvTagsMap = map[string]string{}

	for i := 0; i < valueTypeReflect.NumField(); i++ {
		field := valueTypeReflect.Field(i)

		tag := field.Tag.Get("csv")
		fieldName := field.Name

		if tag == "" || tag == "-" {
			continue
		}

		csvTagsMap[tag] = fieldName
	}
	return csvTagsMap
}

// toMap merge a csv line (which is a slice) into a map, being the keys the header names
// In case Options.UseHeader is false, the keys are the fields positions
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
