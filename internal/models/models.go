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

type HourRow struct{
	Label string
	Slots []int
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

type ScheduleTemplateData struct {
	DayStrings []string
	Days []int
	Slots []int
	HourRows []HourRow
	User User
	Schedule Schedule
	Status string
	Msg *string
}

// type Shift struct {
// 	Id int
// 	StartTime int
// 	Endtime int
// 	Day int
// }