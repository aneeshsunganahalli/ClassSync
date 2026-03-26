package router

import (
	"net/http"

	"github.com/aneeshsunganahalli/ClassSync/internal/api/handlers"
)

func Router() http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.RootHandler)

	mux.HandleFunc("/teachers/", handlers.TeachersHandler)
	mux.HandleFunc("/students", handlers.StudentsHandler)
	mux.HandleFunc("/execs", handlers.ExecsHandler)

	return mux
}