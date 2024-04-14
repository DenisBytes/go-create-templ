package template

import (
	_ "embed"
)

//go:embed files/cmd/migrate/main.go.tmpl
var cmdMigrateTemplate []byte

//go:embed files/cmd/reset/main.go.tmpl
var cmdMigrateReset []byte

func CmdMigrateTemplate() []byte {
	return cmdMigrateTemplate
}

func CmdMigrateReset() []byte {
	return cmdMigrateReset
}