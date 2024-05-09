package template

import (
	_ "embed"
)

//go:embed files/main/gin.go.tmpl
var ginMainTemplate []byte

//go:embed files/handlers/gin/auth.go.tmpl
var ginHandlerAuth []byte

//go:embed files/handlers/gin/home.go.tmpl
var ginHandlerHome []byte

//go:embed files/handlers/gin/app.go.tmpl
var ginHandlerApp []byte

//go:embed files/handlers/gin/middleware.go.tmpl
var ginHandlerMiddleware []byte

//go:embed files/handlers/gin/settings.go.tmpl
var ginHandlerSettings []byte

//go:embed files/handlers/gin/util.go.tmpl
var ginHandlerUtil []byte

type GinTemplates struct{}

func (c GinTemplates) Main() []byte {
	return ginMainTemplate
}

func (c GinTemplates) HandlerAuth() []byte {
	return ginHandlerAuth
}

func (c GinTemplates) HandlerHome() []byte {
	return ginHandlerHome
}

func (c GinTemplates) HandlerApp() []byte{
	return ginHandlerApp
}

func (c GinTemplates) HandlerMiddleware() []byte {
	return ginHandlerMiddleware
}

func (c GinTemplates) HandlerSettings() []byte {
	return ginHandlerSettings
}

func (c GinTemplates) HandlerUtil() []byte {
	return ginHandlerUtil
}