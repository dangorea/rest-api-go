package handlers

import (
	"fmt"
	"net/http"
)

func ExecsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		{
			w.Write([]byte("Hello GET Method on Execs Route"))
		}
	case http.MethodPost:
		{
			fmt.Println("Query:", r.URL.Query())
			fmt.Println("name:", r.URL.Query().Get("name"))

			// Parse form data (necessary for x-www-urlencoded)
			err := r.ParseForm()
			if err != nil {
				return
			}

			fmt.Println("Form from POST methods:", r.Form)

			w.Write([]byte("Hello POST Method on Execs Route"))
		}
	case http.MethodPut:
		{
			w.Write([]byte("Hello POST Method on Execs Route"))
		}
	case http.MethodPatch:
		{
			w.Write([]byte("Hello PATCH Method on Execs Route"))
		}
	case http.MethodDelete:
		{
			w.Write([]byte("Hello DELETE Method on Execs Route"))
		}
	default:
		{
			w.Write([]byte("Hello Execs Route"))
			fmt.Println("Hello Execs Route")
		}
	}
}
