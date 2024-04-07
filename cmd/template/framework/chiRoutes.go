package framework

import (
	_ "embed"
)

//go:embed files/main/chi.go.tmpl
var chiRoutesTemplate []byte

//go:embed files/gitignore.tmpl
var gitIgnoreTemplate []byte

//go:embed files/air.toml.tmpl
var airTomlTemplate []byte

type ChiTemplates struct{}

func (c ChiTemplates) Main() []byte {
	return chiRoutesTemplate
}

func GitIgnoreTemplate() []byte {
	return gitIgnoreTemplate
}

func AirTomlTemplate() []byte {
	return airTomlTemplate
}