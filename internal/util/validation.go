package util

import (
	"errors"
	"net/http"
	"strings"
)

type Source int

const (
	Query Source = iota
	Header
)

func (source Source) String() string {
	switch source {
	case Query:
		return "Query"
	case Header:
		return "Header"
	default:
		return "Invalid"
	}
}

func ValidateParameters(r *http.Request, source Source, fields []string) error {
	var extra []string
	var missing []string
	var input map[string][]string

	switch source {
	case Query:
		input = r.URL.Query()
	case Header:
		input = r.Header
	default:
		return errors.New("invalid parameter source type")
	}

	for _, field := range fields {
		if len(input[field]) == 0 {
			missing = append(missing, field)
		}
		delete(input, field)
	}

	for field := range input {
		extra = append(extra, field)
	}

	parameters := source.String() + " parameters "
	if len(missing) != 0 {
		required := strings.Join(missing, ", ")
		return errors.New("Required " + parameters + required + " must be provided.")
	}

	if len(extra) != 0 {
		additional := strings.Join(extra, ", ")
		return errors.New("Additional " + parameters + additional + " must be removed.")
	}

	return nil
}
