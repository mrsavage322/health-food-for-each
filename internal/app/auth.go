package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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


	err := Config.URLMapDB.SetDB(context.Background(), id, link, app.Cfg.UserID)
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
