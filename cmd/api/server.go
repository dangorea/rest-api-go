package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	mw "rest-api/internal/api/middlewares"
	"rest-api/internal/api/router"
	"rest-api/internal/repository/sqlconnect"

	"github.com/joho/godotenv"
)

type user struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	City string `json:"city"`
}

func main() {

	err := godotenv.Load("cmd/api/.env")

	if err != nil {
		fmt.Println("Error loading .env file")
		fmt.Println(err)
		return
	}

	db, err := sqlconnect.ConnectDb()

	if err != nil {
		fmt.Println("Error-----:", err)
		return
	}

	defer db.Close()

	port := os.Getenv("API_PORT")

	cert := "cert.pem"
	key := "key.pem"

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// rl := mw.NewRateLimiter(5, time.Minute)

	// hppOptions := mw.HPPOptions{
	// 	CheckQuery:                  true,
	// 	CheckBody:                   true,
	// 	CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
	// 	Whitelist:                   []string{"sortBy", "order", "name", "age", "city"},
	// }

	// secureMux := mw.Hpp(hppOptions)(rl.Middleware(mw.Compression(mw.ResponseTimeMiddleware(mw.SecurityHeaders(mw.Cors(mux))))))
	// secureMux := applyMiddlewares(mux, mw.Hpp(hppOptions), mw.Compression, mw.SecurityHeaders, mw.ResponseTimeMiddleware, rl.Middleware, mw.Cors)
	secureMux := mw.SecurityHeaders(router.Router())

	server := &http.Server{
		Addr: fmt.Sprintf(":%s", port),
		// Handler:   middlewares.Cors(mux),
		// Handler: mw.ResponseTimeMiddleware(mw.SecurityHeaders(mw.Cors(mux))),
		Handler:   secureMux,
		TLSConfig: tlsConfig,
	}

	fmt.Println("Server is running on port:", port)

	err = server.ListenAndServeTLS(cert, key)

	if err != nil {
		log.Fatal("Error starting the server", err)
	}
}
