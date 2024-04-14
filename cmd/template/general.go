package template

import (
	_ "embed"
)

//go:embed files/gitignore.tmpl
var gitIgnoreTemplate []byte

//go:embed files/air.toml.tmpl
var airTomlTemplate []byte

//go:embed files/globalEnv.tmpl
var envTemplate []byte

//go:embed files/makefile.tmpl
var makefileTemplate []byte

//go:embed files/tailwind.config.js.tmpl
var tailwindConfigTemplate []byte

func GitIgnoreTemplate() []byte {
	return gitIgnoreTemplate
}

func AirTomlTemplate() []byte {
	return airTomlTemplate
}

func EnvTemplate() []byte {
	return envTemplate
}

func MakefileTemplate() []byte {
	return makefileTemplate
}

func TailwindConfigTemplate() []byte {
	return tailwindConfigTemplate
}