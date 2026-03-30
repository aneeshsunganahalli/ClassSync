package handlers

import (
	"encoding/json"
	"fmt"

	"net/http"

	"strconv"
	"strings"

	"github.com/aneeshsunganahalli/ClassSync/internal/models"
	"github.com/aneeshsunganahalli/ClassSync/internal/repository/sqlconnect"
)

func GetOneTeacherHandler(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")
	fmt.Println(idStr)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err)
		return
	}
	teacher, err := sqlconnect.GetTeacherByID(id)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teacher)
}

func GetTeachersHandler(w http.ResponseWriter, r *http.Request) {
	var teachers []models.Teacher
	teachers, err := sqlconnect.GetAllTeachers(teachers, r)
	if err != nil {
		return
	}

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "Success",
		Count:  len(teachers),
		Data:   teachers,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func AddTeachersHandler(w http.ResponseWriter, r *http.Request) {
	

	var newTeachers []models.Teacher
	err := json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		return
	}

	addedTeachers, err := sqlconnect.AddTeachers(newTeachers)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "Success",
		Count:  len(addedTeachers),
		Data:   addedTeachers,
	}
	json.NewEncoder(w).Encode(response)
}



func UpdateTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/teachers/")
	print(idStr)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Teacher ID ", http.StatusBadRequest)
		return
	}

	var updatedTeacher models.Teacher
	err = json.NewDecoder(r.Body).Decode(&updatedTeacher)
	if err != nil {
		http.Error(w, "Invalid Request Payload", http.StatusBadRequest)
	}

	updatedTeacherFromDb, err := sqlconnect.UpdateTeacher(id, updatedTeacher)
	if err != nil  {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTeacherFromDb)

}



func PatchOneTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Teacher ID ", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Invalid Request Payload ", http.StatusBadRequest)
		return
	}

	existingTeacher, err := sqlconnect.PatchTeacherbyID(id, updates)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingTeacher)
}



func PatchTeachersHandler(w http.ResponseWriter, r *http.Request) {


	var updates []map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

		err = sqlconnect.PatchTeachers(updates)
		if err != nil {
			return
		}

	w.WriteHeader(http.StatusNoContent)
}


func DeleteOneTeacherHandler(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Teacher ID", http.StatusBadRequest)
		return
	}

	err = sqlconnect.DeleteTeacherByID(id)
	if err != nil {
		return
	}

	// w.WriteHeader(http.StatusNoContent)

	w.Header().Set("Content-Type", "application/json")

	response := struct {
		Status string `json:"status"`
		ID     int    `json:"id"`
	}{
		Status: "Teacher deleted successfully",
		ID:     id,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

}



func DeleteTeachersHandler(w http.ResponseWriter, r *http.Request) {
	var ids []int
	if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if len(ids) == 0 {
		http.Error(w, "No teacher IDs provided", http.StatusBadRequest)
		return
	}

	deletedIds, err := sqlconnect.DeleteTeachers(ids)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Status string `json:"status"`
		Count  int    `json:"count"`
		IDs    []int  `json:"ids"`
	}{
		Status: "Teachers deleted successfully",
		Count:  len(deletedIds),
		IDs:    deletedIds,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

