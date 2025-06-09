package system

type _endpoint struct {
	id        int
	path      string
	uriParams int
	methods   int
}

type _method struct {
	id      int
	name    string
	query   int
	headers int
}

type _parameter struct {
	id         int
	name       string
	typ        string
	required   bool
	properties int
}

type _property struct {
	id    int
	name  string
	value string
}
