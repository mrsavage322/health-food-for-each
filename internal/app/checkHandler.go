package app

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log"
	"net/http"
	"time"
)

func BDConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := pgx.Connect(r.Context(), Config.DatabaseAddress)
	if err != nil {
		log.Println("Database connection error:", err)
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer conn.Close(r.Context())

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err = conn.Ping(ctx); err != nil {
		log.Println("Failed to connect to the database:", err)
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
