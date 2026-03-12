package sqlconnect

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"rest-api/internal/models"
	"rest-api/pkg/utils"
	"strconv"
	"time"

	"github.com/go-mail/mail/v2"
)

func GetExecsDbHandler(execs []models.Exec, r *http.Request) ([]models.Exec, error) {
	db, err := ConnectDb()
	if err != nil {
		//http.Error(w, "Database connection error", http.StatusInternalServerError)
		return nil, utils.ErrorHandler(err, "error retrieving data")
	}
	defer db.Close()

	query := "SELECT id, first_name, last_name, email, username, user_created_at, inactive_status, role FROM execs WHERE 1=1"
	var args []interface{}

	query, args = utils.AddFilters(r, query, args)

	query = utils.AddSorting(r, query)

	rows, err := db.Query(query, args...)
	if err != nil {
		fmt.Println(err)
		//http.Error(w, "Database query error", http.StatusInternalServerError)
		return nil, utils.ErrorHandler(err, "error retrieving data")
	}
	defer rows.Close()

	//execList := make([]models.Exec, 0)
	for rows.Next() {
		var exec models.Exec
		err = rows.Scan(&exec.ID, &exec.FirstName, &exec.LastName, &exec.Email, &exec.Username, &exec.UserCreatedAt, &exec.InactiveStatus, &exec.Role)
		if err != nil {
			fmt.Println(err)
			return nil, utils.ErrorHandler(err, "error retrieving data")
		}
		execs = append(execs, exec)
	}
	return execs, nil
}

func GetExecById(id int) (models.Exec, error) {
	db, err := ConnectDb()
	if err != nil {
		return models.Exec{}, utils.ErrorHandler(err, "error retrieving data")
	}
	defer db.Close()

	var exec models.Exec
	err = db.QueryRow("SELECT id, first_name, last_name, email, username, inactive_status, role FROM execs WHERE id = ?", id).Scan(&exec.ID, &exec.FirstName, &exec.LastName, &exec.Email, &exec.Username, &exec.InactiveStatus, &exec.Role)

	if err == sql.ErrNoRows {
		//http.Error(w, "Exec not found", http.StatusNotFound)
		return models.Exec{}, utils.ErrorHandler(err, "error retrieving data")
	} else if err != nil {
		//http.Error(w, "Database connection error", http.StatusInternalServerError)
		return models.Exec{}, utils.ErrorHandler(err, "error retrieving data")
	}
	return exec, nil
}

func AddExecDbHandler(newExecs []models.Exec) ([]models.Exec, error) {
	db, err := ConnectDb()
	if err != nil {
		return nil, utils.ErrorHandler(err, "error adding data")
	}
	defer db.Close()

	stmt, err := db.Prepare(utils.GenerateInsertQuery("execs", models.Exec{}))

	if err != nil {
		return nil, utils.ErrorHandler(err, "error adding data")
	}
	defer stmt.Close()

	addedExecs := make([]models.Exec, len(newExecs))

	for i, newExec := range newExecs {
		newExec.Password, err = utils.HashPassword(newExec.Password)
		if err != nil {
			return nil, utils.ErrorHandler(err, "error adding exec into database")
		}

		structValues := utils.GetStructValues(newExec)
		res, err := stmt.Exec(structValues...)
		if err != nil {
			return nil, utils.ErrorHandler(err, "error adding data")
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			return nil, utils.ErrorHandler(err, "error adding data")
		}
		newExec.ID = int(lastID)
		addedExecs[i] = newExec
	}
	return addedExecs, nil
}

func PatchExecsDbHandler(updates []map[string]interface{}) error {
	db, err := ConnectDb()
	if err != nil {
		return utils.ErrorHandler(err, "error updating data")
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return utils.ErrorHandler(err, "error updating data")
	}

	for _, update := range updates {
		id, ok := update["id"]
		if !ok {
			tx.Rollback()
			return utils.ErrorHandler(err, "Invalid Id")
		}

		var execFromDb models.Exec
		err := db.QueryRow("SELECT id, first_name, last_name, email, username FROM execs WHERE id = ?", id).Scan(&execFromDb.ID, &execFromDb.FirstName, &execFromDb.LastName, &execFromDb.Email, &execFromDb.Username)

		if err != nil {
			tx.Rollback()

			if err == sql.ErrNoRows {
				return utils.ErrorHandler(err, "Exec not found")
			}

			return utils.ErrorHandler(err, "error updating data")
		}

		execVal := reflect.ValueOf(&execFromDb).Elem()
		execType := execVal.Type()

		for k, v := range update {
			if k == "id" {
				continue
			}
			for i := 0; i < execVal.NumField(); i++ {
				field := execType.Field(i)

				if field.Tag.Get("json") == k+",omitempty" {
					fieldValue := execVal.Field(i)
					if execVal.Field(i).CanSet() {
						val := reflect.ValueOf(v)
						if val.Type().ConvertibleTo(fieldValue.Type()) {
							fieldValue.Set(val.Convert(fieldValue.Type()))
						} else {
							tx.Rollback()
							log.Printf("Cannot convert %v to %v", val.Type(), fieldValue.Type())
							return utils.ErrorHandler(err, "error updating data")
						}
					}
					break
				}
			}
		}

		_, err = tx.Exec("UPDATE execs SET first_name = ?, last_name = ?, email = ?, username = ? WHERE id = ?", execFromDb.FirstName, execFromDb.LastName, execFromDb.Email, execFromDb.Username, id)

		if err != nil {
			return utils.ErrorHandler(err, "error updating data")
		}
	}
	err = tx.Commit()
	if err != nil {
		return utils.ErrorHandler(err, "error updating data")
	}
	return nil
}

func PatchExecDbHandler(id int, updates map[string]interface{}) (models.Exec, error) {
	db, err := ConnectDb()
	if err != nil {
		return models.Exec{}, utils.ErrorHandler(err, "error updating data")
	}
	defer db.Close()

	var existingExec models.Exec
	err = db.QueryRow("SELECT id, first_name, last_name, email, username FROM execs WHERE id = ?", id).Scan(&existingExec.ID, &existingExec.FirstName, &existingExec.LastName, &existingExec.Email, &existingExec.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Exec{}, utils.ErrorHandler(err, "Exec not found")
		}

		return models.Exec{}, utils.ErrorHandler(err, "error updating data")
	}

	// Apply updates using reflect package
	execVal := reflect.ValueOf(&existingExec).Elem()
	execType := execVal.Type()

	for k, v := range updates {
		for i := 0; i < execVal.NumField(); i++ {
			field := execType.Field(i)

			if field.Tag.Get("json") == k+",omitempty" {
				if execVal.Field(i).CanSet() {
					fieldValue := execVal.Field(i)
					fieldValue.Set(reflect.ValueOf(v).Convert(execVal.Field(i).Type()))
				}
			}
		}
	}

	_, err = db.Exec("UPDATE execs SET first_name = ?, last_name = ?, email = ?, username = ? WHERE id = ?", existingExec.FirstName, existingExec.LastName, existingExec.Email, existingExec.Username, existingExec.ID)
	if err != nil {
		return models.Exec{}, utils.ErrorHandler(err, "error updating data")
	}
	return existingExec, nil
}

func DeleteExecDbHandler(id int) error {
	db, err := ConnectDb()
	if err != nil {
		return utils.ErrorHandler(err, "error deleting data")
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM execs WHERE id = ?", id)

	if err != nil {
		return utils.ErrorHandler(err, "error deleting data")
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return utils.ErrorHandler(err, "error deleting data")
	}

	if rowsAffected == 0 {
		return utils.ErrorHandler(err, "Exec not found")
	}
	return nil
}

func GetUserByUsername(username string) (*models.Exec, error) {
	db, err := ConnectDb()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Internal server error")
	}
	defer db.Close()

	user := &models.Exec{}
	err = db.QueryRow(`SELECT id, first_name, last_name, email, username, password, inactive_status, role FROM execs WHERE username = ?`, username).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Username, &user.Password, &user.InactiveStatus, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.ErrorHandler(err, "user not found")
		}

		return nil, utils.ErrorHandler(err, "error getting data")
	}
	return user, nil
}

func UpdatePasswordInDb(userId int, currentPassword string, newPassword string) (string, error) {
	db, err := ConnectDb()
	if err != nil {
		return "", utils.ErrorHandler(err, "database connection error")
	}
	defer db.Close()

	var username string
	var userPassword string
	var userRole string

	err = db.QueryRow("SELECT username, password, role FROM execs WHERE id = ?", userId).Scan(&username, &userPassword, &userRole)
	if err != nil {
		return "", utils.ErrorHandler(err, "User not found")
	}

	err = utils.VerifyPassword(currentPassword, userPassword)
	if err != nil {
		return "", utils.ErrorHandler(err, "The password is incorrect")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return "", utils.ErrorHandler(err, "internal error")
	}
	currentTime := time.Now().Format(time.RFC3339)

	_, err = db.Exec("UPDATE execs SET password = ?, password_changed_at = ? WHERE id = ?", hashedPassword, currentTime, userId)
	if err != nil {
		return "", utils.ErrorHandler(err, "failed to update the password")
	}

	token, err := utils.SignToken(userId, username, userRole)
	if err != nil {
		return "", utils.ErrorHandler(err, "Password updated. Could not create token")
	}

	return token, nil
}

func ForgotPasswordDbHandler(emailId string) error {
	db, err := ConnectDb()
	if err != nil {
		return utils.ErrorHandler(err, "internal server error")
	}
	defer db.Close()

	var exec models.Exec
	err = db.QueryRow("SELECT id FROM execs WHERE email = ?", emailId).Scan(&exec.ID)
	if err != nil {
		return utils.ErrorHandler(err, "user not found")
	}

	duration, err := strconv.Atoi(os.Getenv("RESET_TOKEN_EXP_DURATION"))
	if err != nil {
		return utils.ErrorHandler(err, "Failed to send password reset email")
	}
	mins := time.Duration(duration)
	expiry := time.Now().Add(mins * time.Minute).Format(time.RFC3339)

	tokenBytes := make([]byte, 32)
	_, err = rand.Read(tokenBytes)
	if err != nil {
		return utils.ErrorHandler(err, "Failed to send password reset email")
	}

	token := hex.EncodeToString(tokenBytes)
	hashedToken := sha256.Sum256(tokenBytes)

	hashedTokenString := hex.EncodeToString(hashedToken[:])

	_, err = db.Exec("UPDATE execs SET password_reset_token = ?, password_token_expires = ? WHERE id = ?", hashedTokenString, expiry, exec.ID)
	if err != nil {
		return utils.ErrorHandler(err, "Failed to send password reset email")
	}

	// Sent the reset email
	resetUrl := fmt.Sprintf("https://localhost:3002/execs/reset-password/%s", token)
	message := fmt.Sprintf("Forgot your password? Reset your password using the following link: \n%s\nIf you didn't request a password rest, please ignore this email. This link is only valid for %d minutes.", resetUrl, int(mins))

	m := mail.NewMessage()
	m.SetHeader("From", "school_admin@school.con")
	m.SetHeader("To", emailId)
	m.SetHeader("Subject", "Your password reset link")
	m.SetBody("text/html", message)

	d := mail.NewDialer("localhost", 1025, "", "")
	err = d.DialAndSend(m)
	if err != nil {
		return utils.ErrorHandler(err, "Failed to send password reset email")
	}
	return nil
}

func ResetPasswordDbHandler(token, newPassword string) error {
	bytes, err := hex.DecodeString(token)
	if err != nil {
		return utils.ErrorHandler(err, "internal server error")
	}

	hashedToken := sha256.Sum256(bytes)
	hashedTokenString := hex.EncodeToString(hashedToken[:])

	db, err := ConnectDb()
	if err != nil {
		return utils.ErrorHandler(err, "internal server error")
	}
	defer db.Close()

	var user models.Exec

	query := "SELECT id, email FROM execs WHERE password_reset_token = ? AND password_token_expires > ?"
	err = db.QueryRow(query, hashedTokenString, time.Now().Format(time.RFC3339)).Scan(&user.ID, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.ErrorHandler(err, "internal server error")
		}

		return utils.ErrorHandler(err, "invalid or expired reset code")
	}

	// Hash the new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return utils.ErrorHandler(err, "internal server error")
	}

	updateQuery := "UPDATE execs SET password = ?, password_reset_token = NULL, password_token_expires = NULL WHERE id = ?"
	_, err = db.Exec(updateQuery, hashedPassword, user.ID)
	if err != nil {
		return utils.ErrorHandler(err, "internal server error")
	}
	return nil
}
