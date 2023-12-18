package internal

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
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

var request Request

// Авторизация - принимаем на вход json: login - password. Парсим его, хэшируем пароль и проверяем в БД. Ответ - ОК

func getHash(login string, password string) []byte {
	hashPassword := fmt.Sprintf(login + password)
	h := sha256.New()
	h.Write([]byte(hashPassword))
	hash := h.Sum(nil)
	return hash
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashPassword := getHash(request.Login, request.Password)

	pass, err := ConnectionDB.GetAuthData(context.Background(), request.Login, hashPassword)
	if err != nil {
		log.Println(err)
		log.Println(pass)
		//log.Println(hashPassword)
		resp := Response{Result: "Failed!"}
		responseData, err := json.Marshal(resp)
		if err != nil {
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseData)
		return
	} else {
		resp := Response{Result: "Success!"}
		responseData, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		//log.Println(hashPassword)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(responseData)
	}
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashPassword := getHash(request.Login, request.Password)

	er := ConnectionDB.SetAuthData(context.Background(), request.Login, hashPassword)
	if er != nil {
		log.Println("Login already exist")
		resp := Response{Result: "Login already exist!"}
		responseData, err := json.Marshal(resp)
		if err != nil {
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseData)
		return
	}
	resp := Response{Result: "Success!"}
	//log.Println(hashPassword)
	responseData, err := json.Marshal(resp)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}
