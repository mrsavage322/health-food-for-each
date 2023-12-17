package internal

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
)

//Модуль должен принимать на вход json в виде {"login": "xyz", password: "123456789asd"}
//Затем сверяет с БД параметры логина и пароля и отдает ответ Success или Password is incorrect
//
//Если все не ок - ошибка, если ок - успешная авторизация и доступ к таблицам и настройкам

type Request struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Response struct {
	Result string
}

// Авторизация - принимаем на вход json: login - password. Парсим его, хэшируем пароль и проверяем в БД. Ответ - ОК
func SignIn(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var request Request
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//Формируем хэш из логина + пароля
	hashPassword := fmt.Sprintf(request.Login + request.Password)
	h := sha256.New()
	h.Write([]byte(hashPassword))
	dst := h.Sum(nil)

	//Логика обработки с БД
	//if app.Cfg.DatabaseAddr != "" {
	//	err := app.Cfg.URLMapDB.SetDB(context.Background(), id, link, app.Cfg.UserID)
	//	if err != nil {
	//		originalURL, err := app.Cfg.URLMapDB.GetReverse(context.Background(), link, app.Cfg.UserID)
	//		if err != nil {
	//			log.Println(err)
	//			return
	//		}
	//		shortURL := fmt.Sprintf("%s/%s", app.Cfg.BaseURL, originalURL)
	//		resp := Response{Result: shortURL}
	//		responseData, err := json.Marshal(resp)
	//		if err != nil {
	//			return
	//		}
	//
	//		w.Header().Set("Content-Type", "application/json")
	//		w.WriteHeader(http.StatusConflict)
	//		w.Write(responseData)
	//		return
	//	}
	//}

	resp := Response{Result: "Success"}
	responseData, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseData)
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var request Request
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	hashPassword := fmt.Sprintf(request.Login + request.Password)
	h := sha256.New()
	h.Write([]byte(hashPassword))
	dst := h.Sum(nil)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseData)
}
