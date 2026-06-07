package main

import "net/http"

func app(path string, filePathRoot string) http.Handler {
	return http.StripPrefix(path, http.FileServer(http.Dir(filePathRoot)))
}