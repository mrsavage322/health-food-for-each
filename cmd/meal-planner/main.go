package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"health-food-for-each/internal"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var ServerAddress string

func main() {
	//SetFlags()
	//SetConfig()

	//Подключение к БД
	//var once sync.Once
	//once.Do(func() {
	//	Config.URLMapDB = app.NewURLDBStorage(app.Config.DatabaseAddr)
	//})

	r := chi.NewRouter()
	//middleware
	//r.Use(app.LogRequest)
	//r.Use(handler.GzipMiddleware)
	//r.Use(app.AuthMiddleware)

	//Хэндлеры
	//r.Get("/", MainPage)
	r.Get("/auth", internal.SignIn)
	r.Get("/auth", internal.SignUp)
	//r.Get("/info", InfoPage)
	//r.Get("/{id}", Redirect)
	//r.Get("/ping", BDConnection)
	//r.Get("/settings", Settings)
	//r.Get("/refresh", Refresh)
	//r.Get("/my-food", ShowMyFood)
	//r.Post("/my-food", AddToMyFood)
	//r.Delete("/my-food", RemoveMyFood)

	//TODO: вынести в конфиг
	ServerAddress = "http://localhost:8080"
	srv := &http.Server{
		Addr:    ServerAddress,
		Handler: r,
	}

	// Создаем канал для сигналов завершения
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

	log.Printf("Server is listening on %s\n", ServerAddress)

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
	log.Println("Server has stopped.")
}
