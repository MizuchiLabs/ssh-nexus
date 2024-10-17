// Package web embeds the static files for the web app
package web

import (
	"embed"
	"io/fs"

	"github.com/labstack/echo/v5"
)

//go:embed build/*
var build embed.FS
var Static fs.FS

func init() {
	Static = echo.MustSubFS(build, "build")
}
