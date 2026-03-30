package sqlconnect

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/aneeshsunganahalli/ClassSync/internal/models"
)


func GetTeacherByID(id int) (models.Teacher, error) {
	db, err := ConnectDb()
	if err != nil {
		// http.Error(w, "Database connection error", http.StatusInternalServerError)
		return models.Teacher{}, err
	}

	defer db.Close()
	var teacher models.Teacher
	err = db.QueryRow("Select id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
	if err == sql.ErrNoRows {
		// http.Error(w, "Teacher doesn't exist", http.StatusNotFound)
		return models.Teacher{},err
	} else if err != nil {
		fmt.Println(err)
		// http.Error(w, "Database query error", http.StatusInternalServerError)
		return models.Teacher{},err
	}
	return teacher, nil
}


func GetAllTeachers(teachers []models.Teacher, r *http.Request) ([]models.Teacher, error) {
	db, err := ConnectDb()
	if err != nil {
		// http.Error(w, "Database connection error", http.StatusInternalServerError)
		return nil, err
	}

	defer db.Close()
	query := "Select id, first_name, last_name, email, class, subject From teachers Where 1=1"
	var args []interface{}

	query, args = addFilter(r, query, args)

	query = applySorting(r, query)
	rows, err := db.Query(query, args...)
	if err != nil {
		fmt.Println(query)
		fmt.Println(err)
		// http.Error(w, "Database Query Error", http.StatusInternalServerError)
		return nil, err
	}
	defer rows.Close()

	// teacherList := make([]models.Teacher, 0)
	for rows.Next() {
		var teacher models.Teacher
		err := rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
		if err != nil {
			// http.Error(w, "Error scanning database", http.StatusInternalServerError)
			return nil, err
		}
		teachers = append(teachers, teacher)
	}
	if err := rows.Err(); err != nil {
		// http.Error(w, "Database rows error", http.StatusInternalServerError)
		return nil, err
	}
	return teachers, nil
}

func applySorting(r *http.Request, query string) string {
	sortParams := r.URL.Query()["sortby"]
	if len(sortParams) > 0 {
		fmt.Println(sortParams)
		query += " ORDER BY"
		for i, param := range sortParams {
			parts := strings.Split(param, ":")
			if len(parts) != 2 {
				continue
			}
			field, order := parts[0], parts[1]
			if !isValidOrder(order) || !isValidField(field) {
				continue
			}

			if i > 0 {
				query += ","
			}

			query += " " + field + " " + order
		}
	}
	return query
}

func AddTeachers(newTeachers []models.Teacher) ([]models.Teacher, error) {
	db, err := ConnectDb()
	if err != nil {
		// http.Error(w, "Database connection error", http.StatusInternalServerError)
		return nil, err
	}

	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES (?,?,?,?,?)")
	if err != nil {
		// http.Error(w, "Error parsing the SQL Query", http.StatusInternalServerError)
		return nil, err
	}
	defer stmt.Close()

	addedTeachers := make([]models.Teacher, len(newTeachers))
	for i, teacher := range newTeachers {
		res, err := stmt.Exec(teacher.FirstName, teacher.LastName, teacher.Email, teacher.Class, teacher.Subject)
		if err != nil {
			// http.Error(w, "Error inserting into database", http.StatusInternalServerError)
			return nil, err
		}
		lastId, err := res.LastInsertId()
		if err != nil {
			// http.Error(w, "Error getting last insert ID", http.StatusInternalServerError)
			return nil, err

		}
		teacher.ID = int(lastId)
		addedTeachers[i] = teacher
	}
	return addedTeachers, nil
}

func UpdateTeacher(id int, updatedTeacher models.Teacher) (models.Teacher, error) {
	db, err := ConnectDb()
	if err != nil {
		// http.Error(w, "Database connection error", http.StatusInternalServerError)
		return models.Teacher{}, err
	}

	defer db.Close()

	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers where id = ?", id).Scan(&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)
	fmt.Println(existingTeacher)

	if err != nil {
		if err == sql.ErrNoRows {
			// http.Error(w, "Teacher doesn't exist", http.StatusNotFound)
			return models.Teacher{}, err
		}
		fmt.Println(err)
		// http.Error(w, "Database query error", http.StatusInternalServerError)
		return models.Teacher{}, err
	}

	updatedTeacher.ID = existingTeacher.ID

	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?", updatedTeacher.FirstName, updatedTeacher.LastName, updatedTeacher.Email, updatedTeacher.Class, updatedTeacher.Subject, updatedTeacher.ID)

	if err != nil {
		// http.Error(w, "Error Updating Teacher ", http.StatusInternalServerError)
		return models.Teacher{}, err
	}
	return updatedTeacher, nil
}

func PatchTeacherbyID(id int, updates map[string]interface{}) (models.Teacher, error) {
	db, err := ConnectDb()
	if err != nil {
		// http.Error(w, "Database connection error", http.StatusInternalServerError)
		return models.Teacher{}, err
	}

	defer db.Close()

	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers where id = ?", id).Scan(&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)

	if err != nil {
		if err == sql.ErrNoRows {
			// http.Error(w, "Teacher doesn't exist", http.StatusNotFound)
			return models.Teacher{}, err
		}
		fmt.Println(err)
		// http.Error(w, "Database query error", http.StatusInternalServerError)
		return models.Teacher{}, err
	}

	teacherVal := reflect.ValueOf(&existingTeacher).Elem()
	teacherType := teacherVal.Type()
	fmt.Println(teacherType)

	for k, v := range updates {
		for i := 0; i < teacherVal.NumField(); i++ {
			field := teacherType.Field(i)

			if field.Tag.Get("json") == k+",omitempty" {
				if teacherVal.Field(i).CanSet() {

					fieldVal := teacherVal.Field(i)

					fieldVal.Set(reflect.ValueOf(v).Convert(fieldVal.Type()))
				}
			}
		}
	}


	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?", existingTeacher.FirstName, existingTeacher.LastName, existingTeacher.Email, existingTeacher.Class, existingTeacher.Subject, existingTeacher.ID)

	if err != nil {
		// http.Error(w, "Error Updating Teacher ", http.StatusInternalServerError)
		return models.Teacher{}, err
	}
	return existingTeacher, nil
}

func PatchTeachers(updates []map[string]interface{}) error {
	db, err := ConnectDb()
	if err != nil {
		// http.Error(w, "Database connection error", http.StatusInternalServerError)
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		// http.Error(w, "Error starting transaction", http.StatusInternalServerError)
		return err
	}

	for _, update := range updates {
		idStr, ok := update["id"].(float64) // Recieves float64 from JSON
		if !ok {
			tx.Rollback()
			// http.Error(w, "Invalid Teacher ID in update", http.StatusBadRequest)
			return err
		}

		id := int(idStr)

		var teacherFromDb models.Teacher
		err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(&teacherFromDb.ID, &teacherFromDb.FirstName, &teacherFromDb.LastName, &teacherFromDb.Email, &teacherFromDb.Class, &teacherFromDb.Subject)

		if err != nil {
			tx.Rollback()
			if err == sql.ErrNoRows {
				// http.Error(w, "Teacher not found", http.StatusNotFound)
				return err
			}
			// http.Error(w, "Error retrieving teacher", http.StatusInternalServerError)
			return err
		}

		teacherVal := reflect.ValueOf(&teacherFromDb).Elem()
		teacherType := teacherVal.Type()
		fmt.Println(teacherType)

		for k, v := range update {
			if k == "id" {
				continue
			}
			for i := 0; i < teacherVal.NumField(); i++ {
				field := teacherType.Field(i)

				if field.Tag.Get("json") == k+",omitempty" {

					fieldVal := teacherVal.Field(i)
					if teacherVal.Field(i).CanSet() {
						val := reflect.ValueOf(v)
						if val.Type().ConvertibleTo(fieldVal.Type()) {
							fieldVal.Set(val.Convert(fieldVal.Type()))
						} else {
							tx.Rollback()
							log.Printf("Cannot convert %v to %v", val.Type(), fieldVal.Type())
							// http.Error(w, "Invalid field type in update payload", http.StatusBadRequest)
							return err
						}

					}
					break
				}
			}
		}

		_, err = tx.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?", teacherFromDb.FirstName, teacherFromDb.LastName, teacherFromDb.Email, teacherFromDb.Class, teacherFromDb.Subject, teacherFromDb.ID)
		if err != nil {
			tx.Rollback()
			// http.Error(w, "Error udpdating teacher", http.StatusInternalServerError)
			return err
		}

	}

	err = tx.Commit()
	if err != nil {
		// http.Error(w, "Error commiting transaction", http.StatusInternalServerError)
		return err
	}
	return nil
}

func DeleteTeacherByID( id int) error {
	db, err := ConnectDb()
	if err != nil {
		// http.Error(w, "Database connection error", http.StatusInternalServerError)
		return err
	}

	defer db.Close()

	res, err := db.Exec("DELETE FROM teachers WHERE id = ?", id)
	if err != nil {
		// http.Error(w, "Error Deleting Teacher", http.StatusInternalServerError)
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		// http.Error(w, "Error getting delete result", http.StatusInternalServerError)
		return err
	}

	if rows == 0 {
		// http.Error(w, "Teacher doesn't exist", http.StatusNotFound)
		return err
	}
	return nil
}

func DeleteTeachers(ids []int) ([]int, error) {
	db, err := ConnectDb()
	if err != nil {
		// http.Error(w, "Database connection error", http.StatusInternalServerError)
		return nil, err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		// http.Error(w, "Error starting transaction", http.StatusInternalServerError)
		return nil, err
	}

	stmt, err := tx.Prepare("DELETE FROM teachers WHERE id = ?")
	if err != nil {
		tx.Rollback()
		// http.Error(w, "Error preparing delete query", http.StatusInternalServerError)
		return nil, err
	}
	defer stmt.Close()

	deletedIds := []int{}

	for _, id := range ids {
		result, err := stmt.Exec(id)
		if err != nil {
			tx.Rollback()
			// http.Error(w, "Error deleting teacher", http.StatusInternalServerError)
			return nil, err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			tx.Rollback()
			// http.Error(w, "Error checking delete result", http.StatusInternalServerError)
			return nil, err
		}

		if rowsAffected > 0 {
			deletedIds = append(deletedIds, id)
		}

		if rowsAffected < 1 {
			tx.Rollback()
			// http.Error(w, fmt.Sprintf("ID %d not found", id), http.StatusNotFound)
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		// http.Error(w, "Error committing transaction", http.StatusInternalServerError)
		return nil, err
	}

	if len(deletedIds) < 1 {
		// http.Error(w, "IDs do not exist", http.StatusNotFound)
		return nil, err
	}
	return deletedIds, nil
}


func addFilter(r *http.Request, query string, args []interface{}) (string, []interface{}) {
	var params = map[string]string{
		"first_name": "first_name",
		"last_name":  "last_name",
		"email":      "email",
		"class":      "class",
		"subject":    "subject",
	}

	for param, dbField := range params {
		value := r.URL.Query().Get(param)
		if value == "" {
			continue
		}
		query += " AND " + dbField + " = ? "
		args = append(args, value)
	}
	return query, args
}

func isValidOrder(order string) bool {
	return order == "asc" || order == "desc"
}

func isValidField(field string) bool {
	var validString = map[string]bool{
		"first_name": true,
		"last_name":  true,
		"email":      true,
		"class":      true,
		"subject":    true,
	}
	return validString[field]
}