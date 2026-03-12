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
	mux.HandleFunc("POST /execs/{id}/update-password", handlers.UpdatePasswordHandler)

	mux.HandleFunc("POST /execs/login", handlers.LoginHandler)
	mux.HandleFunc("POST /execs/logout", handlers.LogoutHandler)
	mux.HandleFunc("POST /execs/forgot-password", handlers.ForgotPasswordHandler)
	mux.HandleFunc("POST /execs/reset-password/{reset_code}/", handlers.ResetPasswordHandler)

	return mux
}
