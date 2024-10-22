package main

import (
	"net/http"
)

func (app *application) reportError(w http.ResponseWriter, r *http.Request, code int, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(),
		"method", method,
		"uri", uri,
		// "trace", string(debug.Stack()),
	)

	http.Error(w, http.StatusText(code), code)
}
