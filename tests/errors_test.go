package test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	parser "structCsvParser"
)

func TestErrors(t *testing.T) {
	t.Run("Wrong time layout", func(t *testing.T) {
		var reader = strings.NewReader(CsvWithHeader)
		type User struct {
			ID        int       `csv:"id"`
			FirstName string    `csv:"first_name"`
			LastName  string    `csv:"last_name"`
			Username  string    `csv:"username"`
			CreatedAt time.Time `csv:"created_at"`
		}

		options := parser.Options{
			UseHeader:  true,
			TimeLayout: "2006",
		}

		csvParser := getParser(t, reader, options)

		var user User
		err := csvParser.ReadInto(&user)
		assert.Error(t, err)

		var parseError parser.ParseError
		assert.True(t, errors.As(err, &parseError))
		assert.True(t, strings.Contains(parseError.Message, "time.Time"))
	})

	t.Run("Passing anything by value", func(t *testing.T) {
		var reader = strings.NewReader(CsvWithHeader)
		type User struct {
			ID        int       `csv:"id"`
			FirstName string    `csv:"first_name"`
			LastName  string    `csv:"last_name"`
			Username  string    `csv:"username"`
			CreatedAt time.Time `csv:"created_at"`
		}

		options := parser.Options{
			UseHeader:  true,
			TimeLayout: "2006-01-02 15:04:05",
		}

		csvParser := getParser(t, reader, options)

		var user User
		inputs := []interface{}{
			user,
			"string",
			true,
			[]User{user},
		}

		for _, i := range inputs {
			err := csvParser.ReadInto(i)
			assert.Error(t, err)

			var parseError parser.ParseError
			assert.True(t, errors.As(err, &parseError))
			assert.Equal(t, parseError.Message, "value must be a non-nil pointer to a struct")
		}
	})

	t.Run("Passing nil", func(t *testing.T) {
		var reader = strings.NewReader(CsvWithHeader)

		options := parser.Options{
			UseHeader:  true,
			TimeLayout: "2006-01-02 15:04:05",
		}

		csvParser := getParser(t, reader, options)

		err := csvParser.ReadInto(nil)
		assert.Error(t, err)

		var parseError parser.ParseError
		assert.True(t, errors.As(err, &parseError))
		assert.Equal(t, parseError.Message, "value must be a non-nil pointer to a struct")
	})

	t.Run("Passing reference to a slice", func(t *testing.T) {
		var reader = strings.NewReader(CsvWithHeader)

		options := parser.Options{
			UseHeader:  true,
			TimeLayout: "2006-01-02 15:04:05",
		}

		csvParser := getParser(t, reader, options)

		var s []string
		err := csvParser.ReadInto(s)
		assert.Error(t, err)

		var parseError parser.ParseError
		assert.True(t, errors.As(err, &parseError))
		assert.Equal(t, parseError.Message, "value must be a non-nil pointer to a struct")
	})
}
