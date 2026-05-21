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

type HourRow struct{
	Label string
	Slots []int
}

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

	log.Print("in Get Schedule")

	files := []string{
		"templates/base.html",
		"templates/schedule.html",
		"templates/week_view.html",
	}
	
	ts, err := template.ParseFiles(files...)

	if err != nil {
		logAndSendError(w, err, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	mu.Lock()
	schedule, ok := schedules[currentUser.ScheduleId]
	var status string
	if ok {
		status = schedule.Status
	}
	mu.Unlock()

	data := struct {
		Days []int
		Slots []int
		HourRows []HourRow
		Schedule models.Schedule
		Status string
	}{
		Days: makeRange(0, 5), 
		Slots: makeRange(0, 60),
		HourRows: []HourRow{
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
		Schedule: schedule,
		Status: status,
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		logAndSendError(w, err, "Internal Server Error", http.StatusInternalServerError)
	}
}

func postSchedule(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	// jsonData := []byte(r.FormValue("slots"))
	log.Print("In Post")
	log.Print(r.FormValue("slots"))
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

	// mu.Lock()
	// for k, v := range schedules[currentUser.ScheduleId].Slots {
	// 	newSlots[k] = v
	// }
	// mu.Unlock()

	valid := validateSchedule(newSlots)
	if valid {
		log.Print("Valid")
		mu.Lock()
		newSchedule := models.Schedule{
			Id: currentUser.ScheduleId,
			UserId: currentUser.Id,
			Status: "Pending",
			Slots: newSlots,
		}
		schedules[currentUser.ScheduleId] = newSchedule
		mu.Unlock()
	} else {
		log.Print("Not Valid")
	}

	

	// fmt.Fprintf(w, "Posted Schedule")
}

func validateSchedule(slots map[string]bool) bool {
	log.Print(slots)
	for day:=0; day<5; day++ {
		count := 0
		for slot:=0; slot<60; slot++ {
			key := strconv.Itoa(day) + "T" + strconv.Itoa(slot)
			_, ok := slots[key]
			if ok {
				count++
			}
			if (count > 0 && !ok) || (slot == 59 && count > 0) {
				if count > 54 || count < 18 {
					return false
				}
				count = 0
			}
		}
	}
	return true
}

func getApproval(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Approval")
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

	test_shift := make(map[string]bool)
	test_shift["1T4"] = true

	schedules[0] = models.Schedule{Id: 0, UserId: 0, Status: "Draft", Slots: test_shift}
	schedules[1] = models.Schedule{Id: 1, UserId: 1, Status: "Draft", Slots: make(map[string]bool)}
	currentUser = users[0]

	fileServer := http.FileServer(http.Dir("static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /schedule/{id}", getSchedule)
	mux.HandleFunc("POST /schedule/{id}", postSchedule)
	mux.HandleFunc("GET /approval", getApproval)

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}