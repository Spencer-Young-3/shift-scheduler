package main

import (
	"fmt"
	"log"
	"html/template"	
	"encoding/json"
	"net/http"
	"shiftscheduler.youngs3.byu.edu/internal/models"
	"sync"
	"strconv"
)

var (
	users = make(map[int]models.User)
	schedules = make(map[int]models.Schedule)
	currentUser models.User
	mu sync.RWMutex
)


func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome: %v", currentUser.Name)
}

func getSchedule(w http.ResponseWriter, r *http.Request) {


	files := []string{
		"templates/base.html",
		"templates/schedule.html",
		"templates/week_view.html",
		"templates/schedule_form.html",
	}
	
	ts, err := template.ParseFiles(files...)

	if err != nil {
		logAndSendError(w, err, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	mu.Lock()
	schedule, ok := schedules[currentUser.ScheduleId]
	var status string
	var msg *string
	if ok {
		status = schedule.Status
		msg = schedule.Msg
	}
	user := currentUser
	mu.Unlock()

	data := models.ScheduleTemplateData{
		DayStrings: []string{
			"Mon.",
			"Tue.",
			"Wed.",
			"Thu.",
			"Fri.",
		},
		Days: makeRange(0, 5), 
		Slots: makeRange(0, 60),
		HourRows: []models.HourRow{
            {Label: "8am",  Slots: []int{0, 1, 2, 3, 4, 5}},
            {Label: "9am",  Slots: []int{6, 7, 8, 9, 10, 11}},
            {Label: "10am", Slots: []int{12, 13, 14, 15, 16, 17}},
            {Label: "11am", Slots: []int{18, 19, 20, 21, 22, 23}},
            {Label: "12pm", Slots: []int{24, 25, 26, 27, 28, 29}},
            {Label: "1pm",  Slots: []int{30, 31, 32, 33, 34, 35}},
            {Label: "2pm",  Slots: []int{36, 37, 38, 39, 40, 41}},
            {Label: "3pm",  Slots: []int{42, 43, 44, 45, 46, 47}},
            {Label: "4pm",  Slots: []int{48, 49, 50, 51, 52, 53}},
            {Label: "5pm",  Slots: []int{54, 55, 56, 57, 58, 59}},
        },
		User: user,
		Schedule: schedule,
		Status: status,
		Msg: msg,
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		logAndSendError(w, err, "Internal Server Error", http.StatusInternalServerError)
	}
}

func postSchedule(w http.ResponseWriter, r *http.Request) {

	var slots []string
	err := json.Unmarshal([]byte(r.FormValue("slots")), &slots)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	newSlots := make(map[string]bool)
	for i:=0; i < len(slots); i++ {
		newSlots[slots[i]] = true
	}

	valid, msg := validateSchedule(newSlots)
	if valid {
		w.WriteHeader(http.StatusCreated)
		log.Print("Valid")
		mu.Lock()
		newSchedule := models.Schedule{
			Id: currentUser.ScheduleId,
			UserId: currentUser.Id,
			Status: "Pending",
			Slots: newSlots,
			Msg: nil,
		}
		schedules[currentUser.ScheduleId] = newSchedule
		mu.Unlock()
	} else {
		log.Print("Not Valid")
		mu.Lock()
		newSchedule := models.Schedule{
			Id: currentUser.ScheduleId,
			UserId: currentUser.Id,
			Status: "Draft",
			Slots: newSlots,
			Msg: msg,
		}
		schedules[currentUser.ScheduleId] = newSchedule
		mu.Unlock()
	}

	mu.Lock()
	id := currentUser.ScheduleId
	mu.Unlock()

	data := createWeekTemplateData(id)

	files := []string{
		"templates/week_view.html",
		"templates/schedule_form.html",
	}
	
	ts, err := template.ParseFiles(files...)
	if err != nil {
		logAndSendError(w, err, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = ts.ExecuteTemplate(w, "schedule_form", data)
	if err != nil {
		logAndSendError(w, err, "Internal Server Error", http.StatusInternalServerError)
	}

}

func validateSchedule(slots map[string]bool) (bool, *string) {

	overall_count := 0
	for day:=0; day<5; day++ {
		count := 0
		for slot:=0; slot<60; slot++ {
			key := strconv.Itoa(day) + "T" + strconv.Itoa(slot)
			_, ok := slots[key]
			if ok {
				count++
				overall_count++
			}
			if (count > 0 && !ok) || (slot == 59 && count > 0) {
				if count > 54 {
					msg := "Shift longer than 9 hours"
					return false, &msg
				}
				if count < 18 {
					msg := "Shift shorter than 3 hours"
					return false, &msg
				}
				count = 0
			}
		}
	}
	if overall_count < 120 {
		msg := "Less than 20 hours"
		return false, &msg
	}
	if overall_count > 240 {
		msg := "More than 40 hours"
		return false, &msg
	}
	return true, nil
}

func getApprovalList(w http.ResponseWriter, r *http.Request) {

	if currentUser.Role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	files := []string{
		"templates/base.html",
		"templates/approval_list.html",
	}
	
	ts, err := template.ParseFiles(files...)

	if err != nil {
		logAndSendError(w, err, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	mu.Lock()
	data := struct {
		ScheduleListRange []int
		User models.User
	} {
		ScheduleListRange: makeRange(0, len(schedules)),
		User: currentUser,
	}
	mu.Unlock()

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		logAndSendError(w, err, "Internal Server Error", http.StatusInternalServerError)
	}
}

func getApproval(w http.ResponseWriter, r *http.Request) {

	if currentUser.Role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	files := []string{
		"templates/base.html",
		"templates/approval.html",
		"templates/week_view.html",
		"templates/schedule_form.html",
	}
	
	ts, err := template.ParseFiles(files...)

	if err != nil {
		logAndSendError(w, err, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	id, err := strconv.Atoi(r.PathValue("schedule_id"))
	if err != nil {
		logAndSendError(w, err, "ID not valid", http.StatusInternalServerError)
		return
	}

	data := createWeekTemplateData(id)

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		logAndSendError(w, err, "Internal Server Error", http.StatusInternalServerError)
	}
}

func postApproval(w http.ResponseWriter, r *http.Request) {

	if currentUser.Role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	id, err := strconv.Atoi(r.PathValue("schedule_id"))
	if err != nil {
		logAndSendError(w, err, "ID not valid", http.StatusInternalServerError)
		return
	}

	status := r.FormValue("status")
	msg := r.FormValue("msg")

	mu.Lock()
	schedule, ok := schedules[id]
	var userId int
	var slots map[string]bool
	if ok {
		userId = schedule.UserId
		slots = schedule.Slots
	}
	mu.Unlock()

	newSchedule := models.Schedule{
		Id: id,
		UserId: userId,
		Status: status,
		Slots: slots,
		Msg: &msg,
	}


	mu.Lock()
	schedules[id] = newSchedule
	mu.Unlock()



	data := createWeekTemplateData(id)

	ts, err := template.ParseFiles("templates/week_view.html")
	if err != nil {
		logAndSendError(w, err, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "week_view", data)
	if err != nil {
		logAndSendError(w, err, "Internal Server Error", http.StatusInternalServerError)
	}
}

func switchUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("otherId"))
	if err != nil {
		logAndSendError(w, err, "Not a valid ID number", http.StatusInternalServerError)
	}

	mu.Lock()
	currentUser = users[id]
	mu.Unlock()
	w.Header().Add("HX-Refresh", "true")
	w.WriteHeader(http.StatusCreated)
}

func createWeekTemplateData(id int) models.ScheduleTemplateData {
	mu.Lock()
	schedule, ok := schedules[id]
	var status string
	var msg *string
	if ok {
		status = schedule.Status
		msg = schedule.Msg
	}
	user := currentUser
	mu.Unlock()

	data := models.ScheduleTemplateData{
		DayStrings: []string{
			"Mon.",
			"Tue.",
			"Wed.",
			"Thu.",
			"Fri.",
		},
		Days: makeRange(0, 5), 
		Slots: makeRange(0, 60),
		HourRows: []models.HourRow{
            {Label: "8am",  Slots: []int{0, 1, 2, 3, 4, 5}},
            {Label: "9am",  Slots: []int{6, 7, 8, 9, 10, 11}},
            {Label: "10am", Slots: []int{12, 13, 14, 15, 16, 17}},
            {Label: "11am", Slots: []int{18, 19, 20, 21, 22, 23}},
            {Label: "12pm", Slots: []int{24, 25, 26, 27, 28, 29}},
            {Label: "1pm",  Slots: []int{30, 31, 32, 33, 34, 35}},
            {Label: "2pm",  Slots: []int{36, 37, 38, 39, 40, 41}},
            {Label: "3pm",  Slots: []int{42, 43, 44, 45, 46, 47}},
            {Label: "4pm",  Slots: []int{48, 49, 50, 51, 52, 53}},
            {Label: "5pm",  Slots: []int{54, 55, 56, 57, 58, 59}},
        },
		User: user,
		Schedule: schedule,
		Status: status,
		Msg: msg,
	}

	return data
}

func makeRange(start, end int) []int {
	var arr []int
	for i := start; i < end; i++ {
		arr = append(arr, i)
	}
	return arr
}

func logAndSendError(w http.ResponseWriter, err error, msg string, code int) {
	log.Print(err.Error())
	http.Error(w, msg, code)
}

func main() {
	mux := http.NewServeMux()

	users[0] = models.User{Id: 0, Name: "student1", Role: "student", ScheduleId: 0}
	users[1] = models.User{Id: 1, Name: "admin1", Role: "admin", ScheduleId: 1}	

	log.Print("Listening on Port 4000")

	schedules[0] = models.Schedule{Id: 0, UserId: 0, Status: "Draft", Slots: make(map[string]bool), Msg: nil}
	schedules[1] = models.Schedule{Id: 1, UserId: 1, Status: "Draft", Slots: make(map[string]bool), Msg: nil}
	currentUser = users[0]

	fileServer := http.FileServer(http.Dir("static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /schedule", getSchedule)
	mux.HandleFunc("POST /schedule", postSchedule)
	mux.HandleFunc("POST /switch-user", switchUser)
	mux.HandleFunc("GET /approval", getApprovalList)
	mux.HandleFunc("GET /approval/{schedule_id}", getApproval)
	mux.HandleFunc("POST /approval/{schedule_id}", postApproval)

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}