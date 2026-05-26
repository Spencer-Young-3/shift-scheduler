package models

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
