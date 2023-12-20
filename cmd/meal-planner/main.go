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

var dbConnection string

func main() {

	dbConnection = "postgres://postgres:SberBank_2023@localhost:5432/kbgu"
	internal.ConnectionDB = internal.DataBaseConnection(dbConnection)
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
	r.Get("/", internal.MainPage)
	r.Post("/sign_in", internal.SignIn)
	r.Get("/add", internal.AddFood)
	r.Post("/sign_up", internal.SignUp)
	r.Get("/sign_up", internal.SignUp)
	r.Post("/add", internal.AddFood)
	r.Get("/sign_in", internal.SignIn)
	//r.Get("/info", InfoPage)
	//r.Get("/{id}", Redirect)
	//r.Get("/ping", BDConnection)
	//r.Get("/settings", Settings)
	//r.Get("/refresh", Refresh)
	//r.Get("/my-food", ShowMyFood)
	//r.Post("/my-food", AddToMyFood)
	//r.Delete("/my-food", RemoveMyFood)

	//TODO: вынести в конфиг
	internal.ServerAddress = ":8080"
	srv := &http.Server{
		Addr:    internal.ServerAddress,
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

	log.Printf("Server is listening on %s\n", internal.ServerAddress)

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
	log.Println("Server has stopped.")
}
