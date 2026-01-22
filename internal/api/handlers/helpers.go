package handlers

import (
	"errors"
	"reflect"
	"rest-api/pkg/utils"
	"strings"
)

func CheckBlankFields(value interface{}) error {
	for i := 0; i < reflect.ValueOf(value).NumField(); i++ {
		field := reflect.ValueOf(value).Field(i)
		if field.Kind() == reflect.String && field.String() == "" {
			return utils.ErrorHandler(errors.New("all Fields are required"), "All Fields are required")
		}
	}
	return nil
}

func GerFieldNames(model interface{}) []string {
	val := reflect.TypeOf(model)
	fields := []string{}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldToAdd := strings.TrimSuffix(field.Tag.Get("db"), ",omitempty")
		fields = append(fields, fieldToAdd)
	}
	return fields
}
