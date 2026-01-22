package sqlconnect

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDb() (*sql.DB, error) {
	fmt.Println("Trying to connect to MariaDB")

	dbuser := os.Getenv("DB_USER")
	dbpassword := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	dbhost := os.Getenv("DB_HOST")
	dbport := os.Getenv("DB_PORT")

	// connectionString := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbuser, dbpassword, dbhost, dbport, dbname)

	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		// panic(err)
		return nil, err
	}
	// defer db.Close()

	fmt.Println("Connected to MariaDB")
	return db, nil
}
