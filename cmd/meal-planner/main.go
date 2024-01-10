package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"health-food-for-each/internal/app"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	app.SetFlags()
	app.SetConfig()
	app.ConnectionDB = app.DataBaseConnection(app.Config.DatabaseAddress)
	r := chi.NewRouter()

	r.Get("/", app.MainPage)
	r.Post("/sign_in", app.SignIn)
	r.Post("/sign_up", app.SignUp)
	r.Post("/food/add", app.AddFood)
	r.Post("/settings", app.Settings)
	r.Post("/food/dislike", app.AddDislikeFood)

	r.Get("/settings", app.Settings)
	r.Get("/calc/day", app.CalculateDayHandler)
	r.Get("/calc/week", app.CalculateWeekHandler)
	r.Get("/food/dislike", app.ShowDislikeFood)
	r.Get("/ping", app.BDConnection)
	r.Get("/food/show", app.ShowFood)

	r.Delete("/food/delete/dislike", app.DeleteDislikeFoodHandler)
	r.Delete("/food/delete", app.DeleteFoodHandler)

	srv := &http.Server{
		Addr:    app.Config.ServerAddress,
		Handler: r,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit

		log.Println("Server is shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("HTTP server shutdown error: %v", err)
		}
	}()

	log.Printf("Server is listening on %s\n", app.ServerAddress)

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
	log.Println("Server has stopped.")
}
