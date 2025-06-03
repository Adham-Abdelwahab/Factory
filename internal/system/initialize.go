package system

import (
	"Factory/internal/util"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

type route struct {
	endpoint  endpoint
	methods   []method
	validator validator
}

func Initialize(r *chi.Mux) {
	registry := endpoints()
	catalog(r, registry)
	system(r, registry)
}

func endpoints() []route {
	var db = util.Database
	var registry []route
	var endpoint endpoint

	stmt := "SELECT * FROM endpoint"
	if rows, err := db.Query(stmt); err == nil {
		_ = db.ForEach(rows, &endpoint, func() error {
			r := loadMethods(endpoint)
			registry = append(registry, r)
			return nil
		})
	} else {
		logrus.Fatal("failed to fetch system endpoints")
	}

	return registry
}

func loadMethods(e endpoint) route {
	var db = util.Database
	var methods []method
	var method method

	var validations = make(map[string]validation)

	stmt := "SELECT * FROM method WHERE id = $1"
	if rows, err := db.Query(stmt, e.methods); err == nil {
		uriParams := loadParameters(e.uriParams)

		_ = db.ForEach(rows, &method, func() error {
			properties := make(map[int]map[string]string)
			methods = append(methods, method)

			headers := loadParameters(method.headers)
			query := loadParameters(method.parameters)

			loadProperties(properties, uriParams)
			loadProperties(properties, headers)
			loadProperties(properties, query)

			validations[method.method] = validation{
				properties,
				uriParams,
				headers,
				query,
			}

			return nil
		})
	} else {
		logrus.Fatal("error retrieving methods for " + e.path)
	}

	return route{e, methods, validations}
}

func loadParameters(id int) []parameter {
	var parameters []parameter
	var parameter parameter
	var db = util.Database

	stmt := "SELECT * FROM parameter WHERE id = $1"
	if rows, err := db.Query(stmt, id); err == nil {
		_ = db.ForEach(rows, &parameter, func() error {
			parameters = append(parameters, parameter)
			return nil
		})
	} else {
		logrus.Fatalf("failed to load parameters with id %v\n", id)
	}

	return parameters
}

func loadProperties(props map[int]map[string]string, parameters []parameter) {
	var db = util.Database

	stmt := "SELECT * FROM property WHERE id = $1"
	for _, parameter := range parameters {
		if parameter.properties == -1 {
			continue
		}

		var properties map[string]string
		var property property

		if rows, err := db.Query(stmt, parameter.properties); err == nil {
			_ = db.ForEach(rows, &property, func() error {
				properties[property.key] = property.value
				return nil
			})
		} else {
			logrus.Fatalf("failed to load properties with id %v\n", parameter.properties)
		}

		props[parameter.properties] = properties
	}
}
