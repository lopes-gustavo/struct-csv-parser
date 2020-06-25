package test

import (
	"github.com/stretchr/testify/assert"
	"io"
	"structCsvParser"
	"testing"
)

var CsvWithHeader = `id,first_name,last_name,username,created_at
1,"Rob","Pike",rob,"2010-01-27 00:00:00"
2,Ken,Thompson,ken,"2010-01-27 00:00:00"
3,"Gustavo","Lopes","lopes-gustavo","2010-01-27 00:00:00"
`

var CsvWithoutHeader = `1,"Rob","Pike",rob,"2010-01-27 00:00:00"
2,Ken,Thompson,ken,"2010-01-27 00:00:00"
3,"Gustavo","Lopes","lopes-gustavo","2010-01-27 00:00:00"
`

func getParser(t *testing.T, reader io.Reader, options structCsvParser.Options) structCsvParser.Parser {
	csvParser, err := structCsvParser.New(reader, options)
	assert.NoError(t, err)

	return csvParser
}
