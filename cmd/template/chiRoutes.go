package template

import (
	_ "embed"
)

//go:embed files/main/chi.go.tmpl
var chiMainTemplate []byte

//go:embed files/handlers/chi/auth.go.tmpl
var chiHandlerAuth []byte

//go:embed files/handlers/chi/home.go.tmpl
var chiHandlerHome []byte

//go:embed files/handlers/chi/app.go.tmpl
var chiHandlerApp []byte

//go:embed files/handlers/chi/middleware.go.tmpl
var chiHandlerMiddleware []byte

//go:embed files/handlers/chi/settings.go.tmpl
var chiHandlerSettings []byte

//go:embed files/handlers/chi/util.go.tmpl
var chiHandlerUtil []byte

type ChiTemplates struct{}

func (c ChiTemplates) Main() []byte {
	return chiMainTemplate
}

func (c ChiTemplates) HandlerAuth() []byte {
	return chiHandlerAuth
}

func (c ChiTemplates) HandlerHome() []byte {
	return chiHandlerHome
}

func (c ChiTemplates) HandlerApp() []byte{
	return chiHandlerApp
}

func (c ChiTemplates) HandlerMiddleware() []byte {
	return chiHandlerMiddleware
}

func (c ChiTemplates) HandlerSettings() []byte {
	return chiHandlerSettings
}

func (c ChiTemplates) HandlerUtil() []byte {
	return chiHandlerUtil
}