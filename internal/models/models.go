package models

// import (
// 	"time"
// )

type User struct {
	Id int
	Name string
	Role string
	ScheduleId int
}

// func (u User) String() string {
// 	return fmt.Sprintf("")
// }

type Schedule struct {
	Id int
	UserId int
	Status string
	Slots map[string]bool
	Msg *string
}

// type Shift struct {
// 	Id int
// 	StartTime int
// 	Endtime int
// 	Day int
// }