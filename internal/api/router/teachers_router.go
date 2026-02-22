package router

import (
	"net/http"
	"rest-api/internal/api/handlers"
)

func teachersRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// Teacher routes
	mux.HandleFunc("GET /teachers", handlers.GetTeachersHandler)
	mux.HandleFunc("POST /teachers", handlers.PostTeacherHandler)
	mux.HandleFunc("PATCH /teachers", handlers.PatchTeachersHandler)
	mux.HandleFunc("DELETE /teachers", handlers.DeleteTeachersHandler)

	mux.HandleFunc("GET /teachers/{id}", handlers.GetTeacherHandler)
	mux.HandleFunc("PUT /teachers/{id}", handlers.UpdateTeacherHandler)
	mux.HandleFunc("PATCH /teachers/{id}", handlers.PatchTeacherHandler)
	mux.HandleFunc("DELETE /teachers/{id}", handlers.DeleteTeacherHandler)

	mux.HandleFunc("GET /teachers/{id}/students", handlers.GetStudentsByTeacherById)
	mux.HandleFunc("GET /teachers/{id}/studentcount", handlers.GetStudentsCountByTeacherId)

	return mux
}
