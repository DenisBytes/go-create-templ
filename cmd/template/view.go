package template

import (
	_ "embed"
)

//go:embed files/view/util.go.tmpl
var viewUtilTemplate []byte

//go:embed files/view/auth/auth.go.tmpl
var viewAuthTemplate []byte

//go:embed files/view/css/app.css.tmpl
var viewCssTemplate []byte

//go:embed files/view/home/index.go.tmpl
var viewHomeTemplate []byte

//go:embed files/view/layout/layout.go.tmpl
var viewLayoutTemplate []byte

//go:embed files/view/settings/account.go.tmpl
var viewAccountTemplate []byte

//go:embed files/view/ui/navigation.go.tmpl
var viewNavigationTemplate []byte

//go:embed files/view/ui/toast.go.tmpl
var viewToastTemplate []byte

func ViewUtilTemplate() []byte {
	return viewUtilTemplate
}

func ViewAuthTemplate() []byte {
	return viewAuthTemplate
}

func ViewCssTemplate() []byte {
	return viewCssTemplate
}

func ViewHomeTemplate() []byte {
	return viewHomeTemplate
}

func ViewLayoutTemplate() []byte {
	return viewLayoutTemplate
}

func ViewAccountTemplate() []byte {
	return viewAccountTemplate
}

func ViewNavigationTemplate() []byte {
	return viewNavigationTemplate
}

func ViewToastTemplate() []byte {
	return viewToastTemplate
}