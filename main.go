package main

import (
	"fmt"
	"log"
	"html/template"	
	"net/http"
	"shiftscheduler.youngs3.byu.edu/internal/models"
	"sync"
)

var (
	users = make(map[int]models.User)
	schedules = make(map[int]models.Schedule)
	shifts = make(map[int]models.Shift)
	currentUser models.User
	mu sync.RWMutex
)


func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome: %v", currentUser.Name)
}

func getSchedule(w http.ResponseWriter, r *http.Request) {

	
	ts, err := template.ParseFiles("templates/schedule.html")

	if err != nil {
		logAndSendError(w, err, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		logAndSendError(w, err, "Internal Server Error", http.StatusInternalServerError)
	}
}

func postSchedule(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Posted Schedule")
}

func getApproval(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Approval")
}

func logAndSendError(w http.ResponseWriter, err error, msg string, code int) {
	log.Print(err.Error())
	http.Error(w, msg, code)
}

func main() {
	mux := http.NewServeMux()

	users[0] = models.User{Id: 0, Name: "student1", Role: "student"}
	users[1] = models.User{Id: 1, Name: "admin1", Role: "admin"}	

	currentUser = users[0]

	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /schedule/{id}", getSchedule)
	mux.HandleFunc("GET /approval", getApproval)

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}