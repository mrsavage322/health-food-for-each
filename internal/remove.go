package internal

import (
	"encoding/json"
	"net/http"
)

//TODO: Проверка для записи уже имеюгося продукта

type NewFoodData struct {
	Foodname    string `json:"foodname"`
}

var nFD NewFoodData

//TODO: !!!
func DeleteFood(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/sign_in", http.StatusSeeOther)
	} else {
		if r.Method == http.MethodPost {
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&nFD)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
