package template

import (
	_ "embed"
)

//go:embed files/main/httpstd.go.tmpl
var httpstdMainTemplate []byte

//go:embed files/handlers/httpstd/auth.go.tmpl
var httpstdHandlerAuth []byte

//go:embed files/handlers/httpstd/home.go.tmpl
var httpstdHandlerHome []byte

//go:embed files/handlers/httpstd/app.go.tmpl
var httpstdHandlerApp []byte

//go:embed files/handlers/httpstd/middleware.go.tmpl
var httpstdHandlerMiddleware []byte

//go:embed files/handlers/httpstd/settings.go.tmpl
var httpstdHandlerSettings []byte

//go:embed files/handlers/httpstd/util.go.tmpl
var httpstdHandlerUtil []byte

type HttpstdTemplates struct{}

func (c HttpstdTemplates) Main() []byte {
	return httpstdMainTemplate
}

func (c HttpstdTemplates) HandlerAuth() []byte {
	return httpstdHandlerAuth
}

func (c HttpstdTemplates) HandlerHome() []byte {
	return httpstdHandlerHome
}

func (c HttpstdTemplates) HandlerApp() []byte{
	return httpstdHandlerApp
}

func (c HttpstdTemplates) HandlerMiddleware() []byte {
	return httpstdHandlerMiddleware
}

func (c HttpstdTemplates) HandlerSettings() []byte {
	return httpstdHandlerSettings
}

func (c HttpstdTemplates) HandlerUtil() []byte {
	return httpstdHandlerUtil
}