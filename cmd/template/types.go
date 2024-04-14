package template

import (
	_ "embed"
)

//go:embed files/types/user.go.tmpl
var typesUserTemplate []byte

//go:embed files/types/account.go.tmpl
var typesAccountTemplate []byte


func TypesUserTemplate() []byte {
	return typesUserTemplate
}

func TypesAccountTemplate() []byte {
	return typesAccountTemplate
}
