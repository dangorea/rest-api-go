package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	mw "rest-api/internal/api/middlewares"
	"strings"
	"time"
)

type user struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	City string `json:"city"`
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintln(w, "Hello Root Route")

	w.Write([]byte("Hello Root Route"))
	fmt.Println("Hello Root Route")
}

func teachersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello Teachers Route")

	switch r.Method {
	case http.MethodGet:
		{
			fmt.Println(r.URL.Path)
			path := strings.TrimPrefix(r.URL.Path, "/teachers/")
			userID := strings.TrimSuffix(path, "/")

			fmt.Println("The ID is:", userID)

			fmt.Println(r.URL.Query())
			queryParams := r.URL.Query()
			sortBy := queryParams.Get("sortby")
			order := queryParams.Get("order")

			fmt.Printf("SortBy: %v, Sort Order: %v", sortBy, order)

			w.Write([]byte("Hello GET Method on Teachers Route"))
		}
	case http.MethodPost:
		{
			// Parse form data (necessary for x-www-urlencoded)
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Error parsing form", http.StatusBadRequest)
				return
			}

			fmt.Println("Form:", r.Form)

			// Prepare response data
			response := make(map[string]interface{})

			for key, value := range r.Form {
				response[key] = value[0]
			}

			fmt.Println("Processed Response Map:", response)

			// Raw Body
			body, err := io.ReadAll(r.Body)
			if err != nil {
				return
			}
			defer r.Body.Close()

			fmt.Println("Raw Body:", body)
			fmt.Println("Raw Body:", string(body))

			// If you expect json data, then unmarshal it
			var userInstance user
			err = json.Unmarshal(body, &userInstance)

			if err != nil {
				return
			}

			fmt.Println("userInstance", userInstance)

			w.Write([]byte("Hello POST Method on Teachers Route"))
		}
	default:
		{
			w.Write([]byte("Hello Teachers Route"))

		}
	}
}

func studentsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello Students Route"))
	fmt.Println("Hello Students Route")
}

func execsHandler(w http.ResponseWriter, r *http.Request) {
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

func main() {
	port := 3000

	cert := "cert.pem"
	key := "key.pem"

	mux := http.NewServeMux()

	mux.HandleFunc("/", rootHandler)

	mux.HandleFunc("/teachers/", teachersHandler)

	mux.HandleFunc("/students/", studentsHandler)

	mux.HandleFunc("/execs/", execsHandler)

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	rl := mw.NewRateLimiter(5, time.Minute)

	hppOptions := mw.HPPOptions{
		CheckQuery:                  true,
		CheckBody:                   true,
		CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
		Whitelist:                   []string{"sortBy", "order", "name", "age", "city"},
	}

	secureMux := mw.Hpp(hppOptions)(rl.Middleware(mw.Compression(mw.ResponseTimeMiddleware(mw.SecurityHeaders(mw.Cors(mux))))))

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
		// Handler:   middlewares.Cors(mux),
		// Handler: mw.ResponseTimeMiddleware(mw.SecurityHeaders(mw.Cors(mux))),
		Handler:   secureMux,
		TLSConfig: tlsConfig,
	}

	fmt.Println("Server is running on port:", port)

	err := server.ListenAndServeTLS(cert, key)

	if err != nil {
		log.Fatal("Error starting the server", err)
	}
}
