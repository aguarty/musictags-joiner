package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

//runServer - RUN!
func (app *application) runServer() {

	server := http.Server{
		Addr:         app.cfg.Server.Host + ":" + app.cfg.Server.Port,
		Handler:      app.createHTTPHandler(),
		ErrorLog:     app.logger.LogError,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	serverErr := make(chan error, 1)
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		serverErr <- server.ListenAndServe()
	}()

	app.logger.Info("Server started in :" + app.cfg.Server.Port + " port")

	select {
	case err := <-serverErr:
		if err != nil {
			app.logger.Fatal("Could not listen this address: ", err.Error())
		}
	case <-osSignals:
		app.logger.Info("Server is shutting down...")
		ctxServ, cancelServ := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancelServ()
		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctxServ); err != nil {
			app.logger.Fatal("Could not gracefully shutdown the server: ", err.Error())
		}
		app.cancel()
	}
	app.logger.Info("Server stoped")
}

//createHTTPHandler create handler
func (app *application) createHTTPHandler() http.Handler {
	mux := chi.NewMux()

	mux.Use(middleware.Recoverer)
	mux.Use(app.logging())
	mux.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}).Handler)

	mux.Route("/api", func(api chi.Router) {
		api.Route("/v1", func(v1 chi.Router) {
			v1.Use(middleware.SetHeader("Content-Type", "application/json; charset=utf-8;"))
			v1.Route("/genres", func(genres chi.Router) {
				genres.Post("/jointags", app.joiningtags)
				genres.Get("/list", app.genresList)
			})
		})
	})
	return mux
}

//sendResponse send response
func (a *application) sendResponse(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		a.logger.Errorf("couldn't send data to connection: %v", err)
	}
}

//logging - middleware for logging
func (a *application) logging() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			defer func() {
				a.logger.Debugf("code: %d, latency: %v, id: %s, %s, %s", ww.Status(), time.Since(start), r.Method, r.URL.Path, r.RemoteAddr)
			}()
			next.ServeHTTP(ww, r)
		})
	}
}
