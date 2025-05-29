package util

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gorilla/schema"
)

func SafeDecode[T interface{}](request *T, r *http.Request) error {
	var fields = make(map[string]string)
	var issues []string
	var err error

	var input = r.URL.Query()
	d := schema.NewDecoder()
	d.IgnoreUnknownKeys(true)

	types := reflect.TypeOf(*request)
	for i := 0; i < types.NumField(); i++ {
		field := types.Field(i)
		name := strings.ToLower(field.Name)
		fields[name] = field.Type.String()
	}

	for key, values := range input {
		name := strings.ToLower(key)
		typ := fields[name]
		if typ == "" {
			continue
		}

		for _, value := range values {
			switch typ {
			case "int":
				_, err = strconv.Atoi(value)
			case "bool":
				_, err = strconv.ParseBool(value)
			default:
				continue
			}

			if err != nil {
				issue := fmt.Sprintf("Failed to convert (%s=%s) to %s.", key, value, typ)
				issues = append(issues, issue)
				err = nil
			}
		}

		delete(fields, name)
	}

	var missing []string
	for field, typ := range fields {
		missing = append(missing, field+":"+typ)
	}

	if len(missing) != 0 {
		parameters := strings.Join(missing, ", ")
		issue := "Required Query Parameters " + parameters + " must be provided."
		issues = append(issues, issue)
	}

	if len(issues) != 0 {
		message := strings.Join(issues, " ")
		return errors.New(message)
	}

	return d.Decode(request, input)
}
