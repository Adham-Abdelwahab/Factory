package system

type endpoint struct {
	id        int
	path      string
	methods   int
	uriParams int
}

type method struct {
	id         int
	method     string
	parameters int
	headers    int
}

type parameter struct {
	id         int
	name       string
	typ        string
	required   bool
	properties int
}

type property struct {
	id    int
	key   string
	value string
}
