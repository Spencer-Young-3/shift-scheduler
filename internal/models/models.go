package models

// import (
// 	"time"
// )

type User struct {
	Id int
	Name string
	Role string
}

// func (u User) String() string {
// 	return fmt.Sprintf("")
// }

type Schedule struct {
	Id int
	UserId int
	Status string
}

type Shift struct {
	Id int
	StartTime int
	Endtime int
	Day int
}