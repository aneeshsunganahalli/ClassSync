package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/aneeshsunganahalli/ClassSync/internal/models"
)

var (
	teachers = make(map[int]models.Teacher)
	mutex = &sync.Mutex{}
	nextId = 1
)

func init() { // No need to manually called init() in Golang
	teachers[nextId] = models.Teacher{
		ID: nextId,
		FirstName: "John",
		LastName: "Jones",
		Class: "9A",
		Subject: "Math",
	}

	nextId++
	teachers[nextId] = models.Teacher{
		ID: nextId,
		FirstName: "Jack",
		LastName: "Jones",
		Class: "9A",
		Subject: "Math",
	}
}


func TeachersHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
		case http.MethodGet:
			getTeachersHandler(w, r)
		case http.MethodPost:
			fmt.Println("Placeholder")
		case http.MethodDelete:
			fmt.Println("Placeholder")
		case http.MethodPut:
			fmt.Println("Placeholder")
		case http.MethodPatch:
			fmt.Println("Placeholder")
	}
}

func getTeachersHandler(w http.ResponseWriter, r *http.Request) {

	lastName := r.URL.Query().Get("last_name")
	firstName := r.URL.Query().Get("first_name")

	fmt.Println(r.URL.Query())

	// DEBUG: Look at your terminal when you refresh the browser
    fmt.Printf("Filtering for: Last='%s', First='%s'\n", lastName, firstName)

	teacherList := make([]models.Teacher, 0, len(teachers))
	for _, teacher := range teachers {
		if (lastName == "" || lastName == teacher.LastName) && (firstName == "" || teacher.FirstName == firstName) {
			fmt.Printf("Checking %s %s against %s %s\n", teacher.FirstName, teacher.LastName, firstName, lastName)
			teacherList = append(teacherList, teacher)
		}
	}

	response := struct {
		Status string `json:"status"`
		Count int `json:"count"`
		Data []models.Teacher `json:"data"`
	}{
		Status: "Success",
		Count: len(teacherList),
		Data: teacherList,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}