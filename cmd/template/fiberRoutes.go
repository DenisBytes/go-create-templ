package template

import (
	_ "embed"
)

//go:embed files/main/fiber.go.tmpl
var fiberMainTemplate []byte

//go:embed files/handlers/fiber/auth.go.tmpl
var fiberHandlerAuth []byte

//go:embed files/handlers/fiber/home.go.tmpl
var fiberHandlerHome []byte

//go:embed files/handlers/fiber/app.go.tmpl
var fiberHandlerApp []byte

//go:embed files/handlers/fiber/middleware.go.tmpl
var fiberHandlerMiddleware []byte

//go:embed files/handlers/fiber/settings.go.tmpl
var fiberHandlerSettings []byte

//go:embed files/handlers/fiber/util.go.tmpl
var fiberHandlerUtil []byte

type FiberTemplates struct{}

func (c FiberTemplates) Main() []byte {
	return fiberMainTemplate
}

func (c FiberTemplates) HandlerAuth() []byte {
	return fiberHandlerAuth
}

func (c FiberTemplates) HandlerHome() []byte {
	return fiberHandlerHome
}

func (c FiberTemplates) HandlerApp() []byte{
	return fiberHandlerApp
}

func (c FiberTemplates) HandlerMiddleware() []byte {
	return fiberHandlerMiddleware
}

func (c FiberTemplates) HandlerSettings() []byte {
	return fiberHandlerSettings
}

func (c FiberTemplates) HandlerUtil() []byte {
	return fiberHandlerUtil
}