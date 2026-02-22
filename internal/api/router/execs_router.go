package router

import (
	"net/http"
	"rest-api/internal/api/handlers"
)

func execsRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /execs", handlers.GetExecsHandler)
	mux.HandleFunc("POST /execs", handlers.PostExecHandler)
	mux.HandleFunc("PATCH /execs", handlers.PatchExecsHandler)

	mux.HandleFunc("GET /execs/{id}", handlers.GetExecHandler)
	mux.HandleFunc("PATCH /execs/{id}", handlers.PatchExecHandler)
	mux.HandleFunc("DELETE /execs/{id}", handlers.DeleteExecHandler)
	mux.HandleFunc("DELETE /execs/{id}/update-password", handlers.GetExecsHandler)

	mux.HandleFunc("POST /execs/login", handlers.GetExecsHandler)
	mux.HandleFunc("POST /execs/logout", handlers.GetExecsHandler)
	mux.HandleFunc("POST /execs/forgot-password", handlers.GetExecsHandler)
	mux.HandleFunc("POST /execs/reset-password/{reset_code}", handlers.GetExecsHandler)

	return mux
}
