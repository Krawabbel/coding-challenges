package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"url-shortener/internal/models"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /ping", ping)
	mux.HandleFunc("GET /{short}", app.redirect)
	mux.HandleFunc("POST /create", app.create)

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}

type CreateData struct {
	URL string `json:"url"`
}

func (c *CreateData) validate() error {

	u, err := url.Parse(c.URL)
	if err != nil {
		return err
	}

	if u.Scheme == "" {
		u.Scheme = "http"
	}

	c.URL = u.String()

	return nil
}

func (app *application) create(w http.ResponseWriter, r *http.Request) {

	blob, err := io.ReadAll(r.Body)
	if err != nil {
		app.reportError(w, r, http.StatusInternalServerError, err)
		return
	}

	var data CreateData
	if err := json.Unmarshal(blob, &data); err != nil {
		app.reportError(w, r, http.StatusNotAcceptable, err)
		return
	}

	if err := data.validate(); err != nil {
		app.reportError(w, r, http.StatusNotAcceptable, err)
		return
	}

	short := "xyz"

	creator := r.RemoteAddr

	hurl, err := app.urls.Insert(short, data.URL, creator)
	if err != nil {
		app.reportError(w, r, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(hurl); err != nil {
		app.logger.Error(err.Error())
		return
	}

}

func (app *application) redirect(w http.ResponseWriter, r *http.Request) {

	short := fmt.Sprint(r.PathValue("short"))

	hurl, err := app.urls.Get(short)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.reportError(w, r, http.StatusNotFound, err)
		} else {
			app.reportError(w, r, http.StatusInternalServerError, err)
		}
		return
	}

	w.Header().Add("Location", hurl.URL)
	w.WriteHeader(http.StatusMovedPermanently)

	if _, err := fmt.Fprintf(w, "Redirecting to %s", hurl.URL); err != nil {
		app.logger.Error(err.Error())
	}
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// snippets, err := app.snippets.Latest()
	// if err != nil {
	// 	app.serverError(w, r, err)
	// 	return
	// }

	// data := app.newTemplateData(r)
	// data.Snippets = snippets

	// app.render(w, r, http.StatusOK, "home.tmpl", data)

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "coming soon...")
}

func ping(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("OK"))
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.reportError(w, r, http.StatusInternalServerError, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		proto := r.Proto
		method := r.Method
		uri := r.URL.RequestURI()

		app.logger.Info("received request", "ip", ip, "proto", proto, "method", method, "uri", uri)

		next.ServeHTTP(w, r)
	})
}

func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}
