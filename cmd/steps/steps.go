package steps

import "github.com/DenisBytes/go-create-app/cmd/flags"

type StepSchema struct {
	StepName string
	Options  []Item
	Headers  string
	Field    string
}

type Steps struct {
	Steps map[string]StepSchema
}

type Item struct {
	Flag, Title, Desc string
}

func InitSteps(projectType flags.Framework) *Steps {
	steps := &Steps{
		map[string]StepSchema{
			"framework": {
				StepName: "Go Project Framework",
				Options: []Item{
					{
						Title: "Standard-library",
						Desc:  "The built-in Go standard library HTTP package",
					},
					{
						Title: "Chi",
						Desc:  "A lightweight, idiomatic and composable router for building Go HTTP services",
					},
					{
						Title: "Gin",
						Desc:  "Features a martini-like API with performance that is up to 40 times faster thanks to httprouter",
					},
					{
						Title: "Fiber",
						Desc:  "An Express inspired web framework built on top of Fasthttp",
					},
					{
						Title: "Echo",
						Desc:  "High performance, extensible, minimalist Go web framework",
					},
				},
				Headers: "What framework do you want to use in your Go project?",
				Field:   projectType.String(),
			},
		},
	}
	return steps
}
