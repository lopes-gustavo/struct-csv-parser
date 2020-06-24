package main

import (
	"fmt"
	"io"
	_ "net/http/pprof"
	"strings"
	"structCsvParser/parser"
	"time"
)

type User struct {
	ID        int       `csv:"id"`
	FirstName string    `csv:"first_name"`
	LastName  string    `csv:"last_name"`
	Username  string    `csv:"username"`
	CreatedAt time.Time `csv:"created_at"`
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
	r := strings.NewReader(withHeader)
	p := parser.New(r, parser.Options{
		UseHeader:  true,
		TimeLayout: "2006-01-02 15:04:05",
	})

	for {
		var user User
		err := p.ReadInto(&user)

		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("%#v\n", &user)
	}
}
