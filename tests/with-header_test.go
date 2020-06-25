package test

import (
	"database/sql"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	parser "structCsvParser"
)

func TestCsvWithHeader(t *testing.T) {
	t.Run("Read 1 line", func(t *testing.T) {
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
		err := csvParser.ReadInto(&user)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, 1)
		assert.Equal(t, user.FirstName, "Rob")
		assert.Equal(t, user.LastName, "Pike")
		assert.Equal(t, user.Username, "rob")
	})

	t.Run("Read 2 Lines", func(t *testing.T) {
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
		err := csvParser.ReadInto(&user)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, 1)
		assert.Equal(t, user.FirstName, "Rob")
		assert.Equal(t, user.LastName, "Pike")
		assert.Equal(t, user.Username, "rob")

		var user2 User
		err2 := csvParser.ReadInto(&user2)
		assert.NoError(t, err2)
		assert.Equal(t, user2.ID, 2)
		assert.Equal(t, user2.FirstName, "Ken")
		assert.Equal(t, user2.LastName, "Thompson")
		assert.Equal(t, user2.Username, "ken")
	})

	t.Run("Can reuse struct", func(t *testing.T) {
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
		err := csvParser.ReadInto(&user)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, 1)
		assert.Equal(t, user.FirstName, "Rob")
		assert.Equal(t, user.LastName, "Pike")
		assert.Equal(t, user.Username, "rob")

		err = csvParser.ReadInto(&user)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, 2)
		assert.Equal(t, user.FirstName, "Ken")
		assert.Equal(t, user.LastName, "Thompson")
		assert.Equal(t, user.Username, "ken")
	})

	t.Run("Throw io.EOF when file is exhausted", func(t *testing.T) {
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
		err := csvParser.ReadInto(&user)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, 1)

		err = csvParser.ReadInto(&user)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, 2)

		err = csvParser.ReadInto(&user)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, 3)

		err = csvParser.ReadInto(&user)
		assert.Equal(t, err, io.EOF)
	})

	t.Run("EOF continues after file is exhausted", func(t *testing.T) {
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
		err := csvParser.ReadInto(&user)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, 1)

		err = csvParser.ReadInto(&user)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, 2)

		err = csvParser.ReadInto(&user)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, 3)

		err = csvParser.ReadInto(&user)
		assert.Equal(t, err, io.EOF)

		err = csvParser.ReadInto(&user)
		assert.Equal(t, err, io.EOF)
	})

	t.Run("Ignore field if not set", func(t *testing.T) {
		var reader = strings.NewReader(CsvWithHeader)
		type User struct {
			ID        int    `csv:"id"`
			FirstName string `csv:"first_name"`
		}

		options := parser.Options{
			UseHeader:  true,
			TimeLayout: "2006-01-02 15:04:05",
		}

		csvParser := getParser(t, reader, options)

		var user User
		err := csvParser.ReadInto(&user)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, 1)
		assert.Equal(t, user.FirstName, "Rob")
	})

	t.Run("Ignore field without csv tag", func(t *testing.T) {
		var reader = strings.NewReader(CsvWithHeader)
		type User struct {
			ID        int    `csv:"id"`
			FirstName string `csv:"first_name"`
			LastName  string
		}

		options := parser.Options{
			UseHeader:  true,
			TimeLayout: "2006-01-02 15:04:05",
		}

		csvParser := getParser(t, reader, options)

		var user User
		err := csvParser.ReadInto(&user)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, 1)
		assert.Equal(t, user.FirstName, "Rob")
		assert.Equal(t, user.LastName, "")
	})

	t.Run("ignore field if tag not found", func(t *testing.T) {
		var reader = strings.NewReader(CsvWithHeader)
		type User struct {
			ID        int    `csv:"id"`
			FirstName string `csv:"first_name"`
			LastName  string `csv:"not_found"`
		}

		options := parser.Options{
			UseHeader:  true,
			TimeLayout: "2006-01-02 15:04:05",
		}

		csvParser := getParser(t, reader, options)

		var user User
		err := csvParser.ReadInto(&user)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, 1)
		assert.Equal(t, user.FirstName, "Rob")
		assert.Equal(t, user.LastName, "")
	})

	t.Run("should parse using FieldConverters", func(t *testing.T) {
		var reader = strings.NewReader(CsvWithHeader)
		type User struct {
			ID        int          `csv:"id"`
			FirstName string       `csv:"first_name"`
			CreatedAt sql.NullTime `csv:"created_at"`
		}

		createdAtConverter := func(s string) (interface{}, error) {
			layout := "2006-01-02 15:04:05"
			parsedTime, err := time.Parse(layout, s)
			if err != nil {
				return sql.NullTime{}, err
			}
			return sql.NullTime{Time: parsedTime, Valid: true}, nil
		}

		options := parser.Options{
			UseHeader:  true,
			TimeLayout: "2006-01-02 15:04:05",
			FieldConverters: map[string]parser.ConverterFunc{
				"created_at": createdAtConverter,
			},
		}

		csvParser := getParser(t, reader, options)

		var user User
		err := csvParser.ReadInto(&user)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, 1)
		assert.Equal(t, user.FirstName, "Rob")
		assert.True(t, user.CreatedAt.Valid)
	})

	t.Run("should parse using TypeConverters", func(t *testing.T) {
		var reader = strings.NewReader(CsvWithHeader)
		type User struct {
			ID        int    `csv:"id"`
			FirstName string `csv:"first_name"`
			LastName  string `csv:"last_name"`
		}

		stringConverter := func(s string) (interface{}, error) { return strings.ToUpper(s), nil }

		options := parser.Options{
			UseHeader:  true,
			TimeLayout: "2006-01-02 15:04:05",
			TypeConverters: map[string]parser.ConverterFunc{
				"string": stringConverter,
			},
		}

		csvParser := getParser(t, reader, options)

		var user User
		err := csvParser.ReadInto(&user)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, 1)
		assert.Equal(t, user.FirstName, "ROB")
		assert.Equal(t, user.LastName, "PIKE")
	})

	t.Run("FieldConverters have priority over TypeConverters", func(t *testing.T) {
		var reader = strings.NewReader(CsvWithHeader)
		type User struct {
			ID        int    `csv:"id"`
			FirstName string `csv:"first_name"`
			LastName  string `csv:"last_name"`
		}

		stringConverter := func(s string) (interface{}, error) { return strings.ToUpper(s), nil }
		firstNameConverter := func(s string) (interface{}, error) { return "FIRST: " + s, nil }

		options := parser.Options{
			UseHeader:  true,
			TimeLayout: "2006-01-02 15:04:05",
			TypeConverters: map[string]parser.ConverterFunc{
				"string": stringConverter,
			},
			FieldConverters: map[string]parser.ConverterFunc{
				"first_name": firstNameConverter,
			},
		}

		csvParser := getParser(t, reader, options)

		var user User
		err := csvParser.ReadInto(&user)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, 1)
		assert.Equal(t, user.FirstName, "FIRST: Rob")
		assert.Equal(t, user.LastName, "PIKE")
	})
}
