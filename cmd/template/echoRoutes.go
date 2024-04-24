package template

import (
	_ "embed"
)

//go:embed files/main/echo.go.tmpl
var echoMainTemplate []byte

//go:embed files/handlers/echo/auth.go.tmpl
var echoHandlerAuth []byte

//go:embed files/handlers/echo/home.go.tmpl
var echoHandlerHome []byte

//go:embed files/handlers/echo/app.go.tmpl
var echoHandlerApp []byte

//go:embed files/handlers/echo/middleware.go.tmpl
var echoHandlerMiddleware []byte

//go:embed files/handlers/echo/settings.go.tmpl
var echoHandlerSettings []byte

//go:embed files/handlers/echo/util.go.tmpl
var echoHandlerUtil []byte

type EchoTemplates struct{}

func (c EchoTemplates) Main() []byte {
	return echoMainTemplate
}

func (c EchoTemplates) HandlerAuth() []byte {
	return echoHandlerAuth
}

func (c EchoTemplates) HandlerHome() []byte {
	return echoHandlerHome
}

func (c EchoTemplates) HandlerApp() []byte{
	return echoHandlerApp
}

func (c EchoTemplates) HandlerMiddleware() []byte {
	return echoHandlerMiddleware
}

func (c EchoTemplates) HandlerSettings() []byte {
	return echoHandlerSettings
}

func (c EchoTemplates) HandlerUtil() []byte {
	return echoHandlerUtil
}