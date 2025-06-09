package system

import (
	"Factory/internal/util"
	"github.com/go-chi/chi"
)

// JObject is a JSON Object
type JObject map[string]any

// registry stores system metadata in the form of
// endpoints, methods, parameters, and properties
type _registry struct {
	endpoints  map[int]_endpoint             // endpoints  >> [id] --> _endpoint
	methods    map[int]map[string]_method    // methods 	  >> [id] --> [verb] --> _method
	parameters map[int]map[string]_parameter // parameters >> [id] --> [name] --> _parameter
	properties map[int]map[string]string     // properties >> [id] --> map of _property name,value pairs
}

var registry = _registry{
	endpoints:  make(map[int]_endpoint),
	methods:    make(map[int]map[string]_method),
	parameters: make(map[int]map[string]_parameter),
	properties: make(map[int]map[string]string),
}

func Initialize(r *chi.Mux) {
	loadEndpoints()
	loadMethods()
	loadParameters()
	loadProperties()

	catalog(r)
	system(r)
}

func load[T any](table string, processor func(T)) {
	var db = util.Database
	var element T

	stmt := "SELECT * FROM " + table
	if rows, err := db.Query(stmt); err == nil {
		defer rows.Close()
		_ = db.ForEach(rows, &element, func() error {
			processor(element)
			return nil
		})
	} else {
		suffix := table + ": " + err.Error()
		panic("failed to fetch system " + suffix)
	}
}

func loadEndpoints() {
	load[_endpoint]("endpoint", func(e _endpoint) {
		registry.endpoints[e.id] = e
	})
}

func loadMethods() {
	load[_method]("method", func(m _method) {
		if registry.methods[m.id] == nil {
			registry.methods[m.id] = make(map[string]_method)
		}
		registry.methods[m.id][m.name] = m
	})
}

func loadParameters() {
	load[_parameter]("parameter", func(p _parameter) {
		if registry.parameters[p.id] == nil {
			registry.parameters[p.id] = make(map[string]_parameter)
		}
		registry.parameters[p.id][p.name] = p
	})
}

func loadProperties() {
	load[_property]("property", func(p _property) {
		if registry.properties[p.id] == nil {
			registry.properties[p.id] = make(map[string]string)
		}
		registry.properties[p.id][p.name] = p.value
	})
}
