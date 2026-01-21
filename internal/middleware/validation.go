package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationResponse struct {
	Errors []ValidationError `json:"errors"`
}

func ValidateRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" || r.Method == "DELETE" {
			next.ServeHTTP(w, r)
			return
		}

		contentType := r.Header.Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			jsonutil.WriteError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
			return
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			jsonutil.Write(w, http.StatusBadRequest, ValidationResponse{
				Errors: []ValidationError{
					{Field: "body", Message: "Invalid JSON format"},
				},
			})
			return
		}

		ctx := context.WithValue(r.Context(), "validated_body", body)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetValidatedBody(ctx context.Context) map[string]interface{} {
	if body, ok := ctx.Value("validated_body").(map[string]interface{}); ok {
		return body
	}
	return nil
}

func ValidateStruct(v interface{}) []ValidationError {
	var errors []ValidationError
	val := reflect.ValueOf(v)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		fieldName := fieldType.Tag.Get("json")

		if fieldName == "" {
			fieldName = fieldType.Name
		}

		// Required field validation
		required := fieldType.Tag.Get("validate") == "required"
		if required && (field.IsZero() || field.Kind() == reflect.String && field.String() == "") {
			errors = append(errors, ValidationError{
				Field:   fieldName,
				Message: fmt.Sprintf("%s is required", fieldName),
			})
			continue
		}

		// String validations
		if field.Kind() == reflect.String {
			str := field.String()

			// Email validation
			if strings.Contains(fieldType.Tag.Get("validate"), "email") && !isValidEmail(str) {
				errors = append(errors, ValidationError{
					Field:   fieldName,
					Message: "Invalid email format",
				})
			}

			// Min length validation
			if minTag := fieldType.Tag.Get("validate"); strings.Contains(minTag, "min=") {
				minLen := 0
				fmt.Sscanf(minTag, "min=%d", &minLen)
				if len(str) < minLen {
					errors = append(errors, ValidationError{
						Field:   fieldName,
						Message: fmt.Sprintf("%s must be at least %d characters", fieldName, minLen),
					})
				}
			}
		}

		// Number validations
		if field.Kind() == reflect.Float64 || field.Kind() == reflect.Int {
			num := field.Float()

			// Greater than validation
			if gtTag := fieldType.Tag.Get("validate"); strings.Contains(gtTag, "gt=") {
				minVal := 0.0
				fmt.Sscanf(gtTag, "gt=%f", &minVal)
				if num <= minVal {
					errors = append(errors, ValidationError{
						Field:   fieldName,
						Message: fmt.Sprintf("%s must be greater than %f", fieldName, minVal),
					})
				}
			}
		}
	}

	return errors
}

func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}
