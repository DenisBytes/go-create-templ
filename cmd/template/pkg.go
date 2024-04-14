package template

import(
	_ "embed"
)

//go:embed files/pkg/sb/supabase.go.tmpl
var supabaseTemplate []byte

//go:embed files/pkg/util/util.go.tmpl
var pkgUtilTemplate []byte

//go:embed files/pkg/validate/validate.go.tmpl
var validateTemplate []byte

func SupabaseTemplate() []byte {
	return supabaseTemplate
}

func PkgUtilTemplate() []byte {
	return pkgUtilTemplate
}

func ValidateTemplate() []byte {
	return validateTemplate
}