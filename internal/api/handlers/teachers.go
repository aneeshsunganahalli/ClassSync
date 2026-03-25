package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/aneeshsunganahalli/ClassSync/internal/models"
	"github.com/aneeshsunganahalli/ClassSync/internal/repository/sqlconnect"
)

var (
	teachers = make(map[int]models.Teacher)
	mutex    = &sync.Mutex{}
	nextId   = 1
)

func init() { // No need to manually called init() in Golang
	teachers[nextId] = models.Teacher{
		ID:        nextId,
		FirstName: "John",
		LastName:  "Jones",
		Class:     "9A",
		Subject:   "Math",
	}

	nextId++
	teachers[nextId] = models.Teacher{
		ID:        nextId,
		FirstName: "Jack",
		LastName:  "Jones",
		Class:     "9A",
		Subject:   "Math",
	}
}

func TeachersHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		getTeachersHandler(w, r)
	case http.MethodPost:
		addTeachersHandler(w, r)
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
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "Success",
		Count:  len(teacherList),
		Data:   teacherList,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func addTeachersHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sqlconnect.ConnectDb()
	if err != nil {
		return
	}

	defer db.Close()

	var newTeachers []models.Teacher
	err = json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		return
	}

	stmt, err := db.Prepare("INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES (?,?,?,?,?)")
	if err != nil {
		http.Error(w, "Error parsing the SQL Query", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	addedTeachers := make([]models.Teacher, len(newTeachers))
	for i, teacher := range newTeachers {
		res, err := stmt.Exec(teacher.FirstName, teacher.LastName, teacher.Email, teacher.Class, teacher.Subject)
		if err != nil {
			http.Error(w, "Error inserting into database", http.StatusInternalServerError)
			return
		}
		lastId, err := res.LastInsertId()
		if err != nil {
			http.Error(w, "Error getting last insert ID", http.StatusInternalServerError)
		}
		teacher.ID = int(lastId)
		addedTeachers[i] = teacher
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "Success",
		Count: len(addedTeachers),
		Data: addedTeachers,
	}
	json.NewEncoder(w).Encode(response)
}
