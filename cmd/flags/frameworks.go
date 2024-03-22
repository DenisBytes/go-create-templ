package flags

import (
	"fmt"
	"strings"
)

type Framework string

const (
	Chi             Framework = "chi"
	Gin             Framework = "gin"
	Fiber           Framework = "fiber"
	StandardLibrary Framework = "standard-library"
	Echo            Framework = "echo"
)

var AllowedProjectTypes = []string{string(Chi), string(Gin), string(Fiber), string(StandardLibrary), string(Echo)}

func (f Framework) String() string {
	return string(f)
}

func (f *Framework) Type() string {
	return "Framework"
}

func (f *Framework) Set(value string) error {
	for _, project := range AllowedProjectTypes {
		if project == value {
			*f = Framework(value)
			return nil
		}
	}
	return fmt.Errorf("Framework to use. Allowed values: %s", strings.Join(AllowedProjectTypes, ", "))
}
