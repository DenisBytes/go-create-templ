package framework

import (
	_ "embed"
)

//go:embed files/chi.go.tmpl
var chiRoutesTemplate []byte

type ChiTemplates struct{}

func (c ChiTemplates) Routes() []byte {
	return chiRoutesTemplate
}