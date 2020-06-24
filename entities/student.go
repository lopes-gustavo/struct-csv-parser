package entities

import "time"

//first_name,last_name,username
type Student struct {
	ID        int       `csv:"id"`
	FirstName string    `csv:"first_name"`
	LastName  string    `csv:"last_name"`
	Username  string    `csv:"username"`
	CreatedAt time.Time `csv:"id"`
}

func NewStudent(csvLine []string) Student {
	student := Student{
		FirstName: csvLine[0],
		LastName:  csvLine[1],
		Username:  csvLine[2],
	}

	return student
}
