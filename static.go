package main

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed static/*
var static embed.FS
var staticFS, _ = fs.Sub(static, "static")

var StaticFileServer = http.FileServer(http.FS(staticFS))
