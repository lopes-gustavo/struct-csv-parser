package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strings"
	"structCsvParser/entities"
	"structCsvParser/parser"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	in := `ID,FirstName,LastName,Username,CreatedAt
1,"Rob","Pike",rob,"2010-01-27 00:00:00"
2,Ken,Thompson,ken,"2010-01-27 00:00:00"
3,"Robert","Griesemer","gri","2010-01-27 00:00:00"
`
	r := strings.NewReader(in)
	p := parser.New(r, parser.Options{
		UseHeader: true,
	})

	for {
		var student entities.Student
		err := p.Read(&student)

		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(student)
	}
}
