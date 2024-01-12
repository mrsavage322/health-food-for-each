package app

import (
	"encoding/json"
	"net/http"
)

type NewFoodData struct {
	Foodname string `json:"foodname"`
}

var nFD NewFoodData
var removingDislikeFood string

// Хэндлер удаляет получаемый продукт пользователя
func DeleteFoodHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "You are not auth", http.StatusBadRequest)
		return
	}

	if r.Method != http.MethodDelete {
		http.Error(w, "Use DELETE method for this endpoint", http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&nFD)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ok := ConnectionDB.DeleteFood(r.Context(), Request.Login, nFD.Foodname)
	if !ok {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("URLs deleted"))

}

// Хэндлер удаляет получаемый продукт, который пользователь добавил в исключение
func DeleteDislikeFoodHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "You are not auth", http.StatusBadRequest)
		return
	}

	if r.Method != http.MethodDelete {
		http.Error(w, "Use DELETE method for this endpoint", http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&removingDislikeFood)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ok := ConnectionDB.DeleteDislikeFood(r.Context(), Request.Login, removingDislikeFood)
	if !ok {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("URLs deleted"))

}
