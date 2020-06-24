package main

import (
	"fmt"
	"io"
	"log"
	"strings"
	"structCsvParser/parser"
	"time"
)

type UserWithHeader struct {
	ID        int       `csv:"id"`
	FirstName string    `csv:"first_name"`
	LastName  string    `csv:"last_name"`
	Username  string    `csv:"username"`
	CreatedAt time.Time `csv:"created_at"`
}

type UserWithoutHeader struct {
	ID        int       `csv:"0"`
	FirstName string    `csv:"1"`
	LastName  string    `csv:"2"`
	Username  string    `csv:"3"`
	CreatedAt time.Time `csv:"4"`
}

var withHeader = `id,first_name,last_name,username,created_at
1,"Rob","Pike",rob,"2010-01-27 00:00:00"
2,Ken,Thompson,ken,"2010-01-27 00:00:00"
3,"Robert","Griesemer","gri","2010-01-27 00:00:00"
`

var withoutHeader = `1,"Rob","Pike",rob,"2010-01-27 00:00:00"
2,Ken,Thompson,ken,"2010-01-27 00:00:00"
3,"Robert","Griesemer","gri","2010-01-27 00:00:00"
`

func main() {
	runWithHeader()
	runWithoutHeader()
}

func runWithHeader() {
	reader := strings.NewReader(withHeader)
	options := parser.Options{
		UseHeader:  true,
		TimeLayout: "2006-01-02 15:04:05",
	}
	csvParser, err := parser.New(reader, options)
	if err != nil {
		log.Fatal(err)
	}

	for {
		var user UserWithHeader
		err := csvParser.ReadInto(&user)

		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("%#v\n", &user)
	}
}

func runWithoutHeader() {
	reader := strings.NewReader(withoutHeader)
	options := parser.Options{
		UseHeader:  false,
		TimeLayout: "2006-01-02 15:04:05",
	}
	csvParser, err := parser.New(reader, options)
	if err != nil {
		log.Fatal(err)
	}

	for {
		var user UserWithoutHeader
		err := csvParser.ReadInto(&user)

		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("%#v\n", &user)
	}
}
