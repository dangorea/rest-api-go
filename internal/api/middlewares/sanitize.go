package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"rest-api/pkg/utils"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

func XSSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Sanitize the URL path
		sanitizedPath, err := clean(r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println("Original Path:", r.URL.Path)
		fmt.Println("Sanitized Path:", sanitizedPath)

		// Sanitize query params
		params := r.URL.Query()
		sanitizedQuery := make(map[string][]string)
		for key, values := range params {
			sanitizedKey, err := clean(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			var sanitizedValues []string
			for _, value := range values {
				cleanValue, err := clean(value)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				sanitizedValues = append(sanitizedValues, cleanValue.(string))
			}
			sanitizedQuery[sanitizedKey.(string)] = sanitizedValues

			fmt.Printf("Original Query %s: %s\n", key, strings.Join(values, ", "))
			fmt.Printf("Sanitized Query %s: %s\n", sanitizedKey, strings.Join(sanitizedValues, ", "))
		}

		r.URL.Path = sanitizedPath.(string)
		r.URL.RawQuery = url.Values(sanitizedQuery).Encode()

		fmt.Println("Updated URL:", r.URL.String())

		// Sanitize request body
		if r.Header.Get("Content-Type") == "application/json" {
			if r.Body != nil {
				bodyBytes, err := io.ReadAll(r.Body)
				if err != nil {

					http.Error(w, utils.ErrorHandler(err, "Error reading request body.").Error(), http.StatusBadRequest)
					return
				}

				bodyString := strings.TrimSpace(string(bodyBytes))
				fmt.Println("Original body:", bodyString)

				// Reset the request body
				r.Body = io.NopCloser(bytes.NewReader([]byte(bodyString)))

				if len(bodyString) > 0 {
					var inputData interface{}
					err = json.NewDecoder(bytes.NewReader([]byte(bodyString))).Decode(&inputData)
					if err != nil {

						http.Error(w, utils.ErrorHandler(err, "Invalid JSON body.").Error(), http.StatusBadRequest)
						return
					}
					fmt.Println("Original JSON data:", inputData)

					// Sanitize the JSON body
					sanitedData, err := clean(inputData)
					if err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}

					fmt.Println("Sanitized JSON data:", sanitedData)

					// Marshal the sanitized data back to the body:
					sanitizedBody, err := json.Marshal(sanitedData)
					if err != nil {

						http.Error(w, utils.ErrorHandler(err, "Error sanitizing body").Error(), http.StatusBadRequest)
						return
					}

					r.Body = io.NopCloser(bytes.NewReader(sanitizedBody))
					fmt.Println("Sanitized JSON data:", string(sanitizedBody))

				} else {
					log.Println("Request body is empty")
				}
			} else {
				log.Println("No body in the request.")
			}
		} else if r.Header.Get("Content-Type") != "" {
			log.Printf("Received request with unsupported Content-Type: %s. Expected application/json.\n", r.Header.Get("Content-Type"))
			http.Error(w, "Unsupported Content-Type. Please use application/json.", http.StatusUnsupportedMediaType)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Clean sanitizes input data to prevent XSS attacks
func clean(data interface{}) (interface{}, error) {
	switch v := data.(type) {
	case map[string]interface{}:
		{
			for key, value := range v {
				v[key] = sanitizeValue(value)
			}
		}
		return data, nil
	case []interface{}:
		for index, value := range v {
			v[index] = sanitizeValue(value)
		}
		return v, nil
	case string:
		return sanitizeString(v), nil
	default:
		// Error
		return nil, utils.ErrorHandler(fmt.Errorf("unsupported type: %T", data), fmt.Sprintf("unsupported type: %T", data))
	}
}

func sanitizeValue(data interface{}) interface{} {
	switch v := data.(type) {
	case string:
		return sanitizeString(v)
	case map[string]interface{}:
		{
			for key, value := range v {
				v[key] = sanitizeValue(value)
			}
		}
		return v
	case []interface{}:
		for index, value := range v {
			v[index] = sanitizeValue(value)
		}
		return v
	default:
		return v // return v as it is unsupported
	}
}

func sanitizeString(value string) string {
	return bluemonday.UGCPolicy().Sanitize(value)
}
