package gen

import (
	"strings"

	"github.com/hofstadter-io/hof/schema/gen"

	"hof.io/docs/example/schema"
)

// Generator definition
Generator: gen.Generator & {

	// User inputs to this generator
	// -----------------------------

	// The server design conforming to the server schema
	Server: schema.Server

	// Datamodel for types in our server
	Datamodel: schema.Datamodel

	// Base output directory, defaults to current
	Outdir: string | *"./"

	// Required fields for hof
	// ------------------------

	// In is passed to every template
	In: {
		SERVER:    Server
		DM:        Datamodel
		Resources: (schema.#DatamodelToResources & {"Datamodel": DM}).Resources
	}

	// Actual files generated by hof, combined into a single list
	Out: [...gen.File] & _All

	_All: [
		for _, F in _OnceFiles {F},
		for _, F in _RouteFiles {F},
		for _, F in _TypeFiles {F},
		for _, F in _ResourceFiles {F},
	]

	// Note, we can omit Templates, Partials, and Statics
	// since the default values are sufficient for us

	// Internal fields for mapping Input to templates
	// ----------------------------------------------

	// Files that are generated once per server
	_OnceFiles: [...gen.#File] & [
			{
			TemplatePath: "go.mod"
			Filepath:     "go.mod"
		},
		{
			TemplatePath: "server.go"
			Filepath:     "server.go"
		},
		{
			TemplatePath: "router.go"
			Filepath:     "router.go"
		},
		{
			TemplatePath: "middleware.go"
			Filepath:     "middleware.go"
		},
		// a conditional file
		if Server.Auth != _|_ {
			TemplatePath: "auth.go"
			Filepath:     "auth.go"
		},
		{
			TemplatePath: "index.html"
			Filepath:     "client/index.html"
		},
	]

	// Routes, we create a file per route in the Server
	_RouteFiles: [...gen.File] & [
			for _, R in Server.Routes {
			In: {
				ROUTE: {
					R
				}
			}
			TemplatePath: "route.go"
			Filepath:     "routes/\(In.ROUTE.Name).go"
		},
	]

	_TypeFiles: [...gen.File] & [
			for n, M in Datamodel.Models if n != "$hof" {
			In: {
				TYPE: {
					M
					OrderedFields: [ for _, F in M.Fields {F}]
				}
			}
			TemplatePath: "type.go"
			Filepath:     "types/\(In.TYPE.Name).go"
		},
	]

	_ResourceFiles: [...gen.File] & [
			for _, R in In.Resources {
			In: {
				RESOURCE: R
			}
			TemplatePath: "resource.go"
			Filepath:     "resources/\(In.RESOURCE.Name).go"
		},
		// HTML content
		for _, R in In.Resources {
			In: {
				RESOURCE: R
			}
			TemplatePath: "resource.html"
			Filepath:     "client/\(strings.ToLower(In.RESOURCE.Name)).html"
		},
		// HTML content
		for _, R in In.Resources {
			In: {
				RESOURCE: R
			}
			TemplatePath: "resource.js"
			Filepath:     "client/\(strings.ToLower(In.RESOURCE.Name)).js"
		},
	]

	// We'll see how to handle nested or sub-routes in the "full-example" section

	...
}
