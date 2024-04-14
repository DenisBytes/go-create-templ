package template

import (
	_ "embed"
)

//go:embed files/db/db.go.tmpl
var dbTemplate []byte

//go:embed files/db/query.go.tmpl
var dbQueryTemplate []byte

func DbTemplate() []byte {
	return dbTemplate
}

func DbQueryTemplate() []byte {
	return dbQueryTemplate
}