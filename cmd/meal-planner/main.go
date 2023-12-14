package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Авторизация - принимаем на вход json: login - password. Парсим его, хэшируем пароль и записываем в БД. Ответ - ОК
func AuthPage(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req Request
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	link := req.URL
	shortURL := fmt.Sprintf("%s/%s", Config.BaseURL, id)

	if app.Cfg.DatabaseAddr != "" {
		err := app.Cfg.URLMapDB.SetDB(context.Background(), id, link, app.Cfg.UserID)
		if err != nil {
			originalURL, err := app.Cfg.URLMapDB.GetReverse(context.Background(), link, app.Cfg.UserID)
			if err != nil {
				log.Println(err)
				return
			}
			shortURL := fmt.Sprintf("%s/%s", app.Cfg.BaseURL, originalURL)
			resp := Response{Result: shortURL}
			responseData, err := json.Marshal(resp)
			if err != nil {
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			w.Write(responseData)
			return
		}
	} else {
		app.Cfg.URLMap.Set(id, link)
	}

	resp := Response{Result: shortURL}
	responseData, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseData)
}

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}

func main() {
	SetFlags()
	SetConfig()

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
	r.Get("/", MainPage)
	r.Get("/auth", AuthPage)
	r.Get("/info", InfoPage)
	r.Get("/{id}", Redirect)
	r.Get("/ping", BDConnection)
	r.Get("/settings", Settings)
	r.Get("/refresh", Refresh)
	r.Get("/my-food", ShowMyFood)
	r.Post("/my-food", AddToMyFood)
	r.Delete("/my-food", RemoveMyFood)

	srv := &http.Server{
		Addr:    Config.ServerAddr,
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

	log.Printf("Server is listening on %s\n", Config.ServerAddr)

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
	log.Println("Server has stopped.")
}
