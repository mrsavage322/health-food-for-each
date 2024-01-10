package internal

import (
	"context"
	"encoding/json"
	"net/http"
)

//TODO: Проверка для записи уже имеюгося продукта

type NewFoodData struct {
	Foodname string `json:"foodname"`
}

var nFD NewFoodData

// TODO: !!!
func DeleteFoodHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/sign_in", http.StatusSeeOther)
	} else {
		if r.Method == http.MethodDelete {
			err := json.NewDecoder(r.Body).Decode(&nFD)
			if err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}
			er := ConnectionDB.DeleteFood(context.Background(), request.Login, nFD.Foodname)
			if !er {
				http.Error(w, "Invalid input data", http.StatusBadRequest)
			} else {
				w.WriteHeader(http.StatusAccepted)
				w.Write([]byte("URLs deleted"))
			}
		}
	}
}
